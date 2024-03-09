import argparse
import warnings

import pandas as pd
import time
import pandas_ta as pd_ta
import psycopg2
import numpy as np

from alice_blue import *
from ta.momentum import RSIIndicator
from ta.trend import EMAIndicator, ADXIndicator, SMAIndicator, PSARIndicator

from config.config import *
from ta.volatility import BollingerBands, AverageTrueRange
from notification.desktop_notification import DesktopNotification
from setup.logger import log
from utils import utils
from utils.utils import save_df, get_candle_data_from_db_for_instrument, delete_all_data_from_db_table

ADX_INDICATOR_WINDOW = 14
RSI_INDICATOR_WINDOW = 14
FASTER_EMA_WINDOW = 10
SLOWER_EMA_WINDOW = 20
SMA_TREND_INDICATOR_WINDOW = 20
EMA_TREND_INDICATOR_WINDOW = 60
RSI_THRESHOLD = 55.0
ADX_THRESHOLD = 20.0
ADX_SLOPE_THRESHOLD = 0.0
SLOPE_PERIOD = 2
FASTER_EMA_WINDOW_COL = f'ema_{FASTER_EMA_WINDOW}'
SLOWER_EMA_WINDOW_COL = f'ema_{SLOWER_EMA_WINDOW}'
EMA_TREND_INDICATOR_COL = f'ema_{EMA_TREND_INDICATOR_WINDOW}'
SMA_TREND_INDICATOR_COL = f'sma_{SMA_TREND_INDICATOR_WINDOW}'
RSI_INDICATOR_COL = f'rsi_{RSI_INDICATOR_WINDOW}'

LONG = 'long'
SHORT = 'short'
LONG_SHORT = 'long_short'
STRATEGY_NAME = 'strategy_name'
STRATEGY_TYPE = 'strategy_type'

desktop_notification = DesktopNotification()
connection = psycopg2.connect(host=DB_HOST, dbname=DB_NAME, user=DB_USER, password=DB_PASSWORD)
# pd.set_option('mode.chained_assignment', 'raise')
warnings.filterwarnings("ignore", category=RuntimeWarning)

supertrend_7_3_strategy = {STRATEGY_NAME: 'supertrend_7_3', STRATEGY_TYPE: LONG_SHORT}
adx_ema_strategy = {STRATEGY_NAME: 'adx_ema', STRATEGY_TYPE: LONG_SHORT}
adx_ema_basic_strategy = {STRATEGY_NAME: 'adx_ema_basic', STRATEGY_TYPE: LONG_SHORT}

ichimoku_chikou_span_strong_strategy = {STRATEGY_NAME: 'ichimoku_chikou_span_strong', STRATEGY_TYPE: LONG_SHORT}
ichimoku_chikou_span_strongest_strategy = {STRATEGY_NAME: 'ichimoku_chikou_span_strongest', STRATEGY_TYPE: LONG_SHORT}
STRATEGY_LIST = [supertrend_7_3_strategy, adx_ema_strategy, adx_ema_basic_strategy]
ICHIMOKU_STRATEGY_LIST = [ichimoku_chikou_span_strategy, ichimoku_chikou_span_strong_strategy, ichimoku_chikou_span_strongest_strategy]


class AliceBlueDailySignals:
    def __init__(self, alice_blue, limit_rows_flag):
        self._alice_blue = alice_blue
        self.limit_rows_flag = limit_rows_flag
        self._token_mapping = {}

    @staticmethod
    def get_all_instruments():
        cursor = connection.cursor()
        cursor.execute("SELECT * FROM instruments WHERE exchange = 'NSE'")
        instruments = cursor.fetchall()
        return instruments

    @staticmethod
    def get_all_instruments_with_filter():
        cursor = connection.cursor()
        cursor.execute(""" SELECT * FROM instruments
                            WHERE symbol NOT LIKE '% %'
                            AND exchange = 'NSE'
                            """)
        instruments = cursor.fetchall()
        return instruments

    def start_alice_blue_daily_signals_generation(self):
        log.info(':thumbs_up: Starting Alice Blue Daily Signal Generation.')
        delete_all_data_from_db_table(connection, 'alpha_signals_daily')
        all_instrument_with_filter = self.get_all_instruments_with_filter()
        total_number_of_instruments = len(all_instrument_with_filter)
        log.info(f"Generating Daily Signals for {total_number_of_instruments} instruments.")

        for idx, instrument in enumerate(all_instrument_with_filter):
            log.info(f"[bold blue]Running for instrument number: {idx+1} of total {total_number_of_instruments} instruments. instrument: {instrument}[/]")
            log.info(f"[bold magenta] Percentage Complete: {round(100 * (idx + 1)/total_number_of_instruments, 5) }%[/]")
            exchange = instrument[0]
            exchange_code = utils.exchange_codes.get(exchange)
            token = instrument[1]
            symbol = instrument[2]
            log.info(f"Calculating the indicator values for instrument: {instrument}")
            instrument_candles_df = get_candle_data_from_db_for_instrument(connection, 'candles_1d', exchange_code, token, self.limit_rows_flag)

            if instrument_candles_df.empty or instrument_candles_df.shape[0] < 52:
                continue

            instrument_candles_df['exchange'] = exchange
            instrument_candles_df['symbol'] = symbol
            instrument_candles_df = supertrend_indicator_values_df(instrument_candles_df)
            instrument_candles_df = bollinger_bands_indicator_values_df(instrument_candles_df)
            instrument_candles_df = average_true_range_indicator_values_df(instrument_candles_df)
            instrument_candles_df = rsi_indicator_values_df(instrument_candles_df, 14)
            instrument_candles_df = ema_indicator_values_df(instrument_candles_df, [10, 20, 60])
            instrument_candles_df = adx_indicator_values_df(instrument_candles_df, 14)
            instrument_candles_df = sma_indicator_values_df(instrument_candles_df, 20)
            instrument_candles_df.ta.ichimoku(append=True)
            # instrument_candles_df.ta.psar(append=True)
            instrument_candles_df = slope_calc(instrument_candles_df, RSI_INDICATOR_COL, SLOPE_PERIOD)
            instrument_candles_df = slope_calc(instrument_candles_df, 'adx_pos', SLOPE_PERIOD)

            instrument_candles_df = supertrend_7_3_alpha_signals(instrument_candles_df)
            instrument_candles_df = adx_ema_alpha_signals(instrument_candles_df)
            instrument_candles_df = adx_ema_basic_alpha_signals(instrument_candles_df)
            instrument_candles_df = ichimoku_chikou_span_alpha_signals_basic(instrument_candles_df)
            instrument_candles_df = ichimoku_chikou_span_alpha_signals_20_days_high_low_crossover(instrument_candles_df)
            instrument_candles_df = ichimoku_chikou_span_alpha_signals(instrument_candles_df)
            instrument_candles_df = ichimoku_chikou_span_strong_alpha_signals(instrument_candles_df)
            instrument_candles_df = ichimoku_chikou_span_strongest_alpha_signals(instrument_candles_df)
            instrument_candles_df['timestamp'] = instrument_candles_df.index
            save_df(connection, instrument_candles_df, 'alpha_signals_daily', 'IGNORE')


def supertrend_indicator_values_df(df):
    if df.empty:
        return df

    try:
        supertrend_df = pd_ta.supertrend(df["high"], df['low'], df['close'], period=7, multiplier=3)
        supertrend_df.columns = ['supertrend', 'supertrend_direction', 'supert_long', 'supert_short']
        df = df.join(supertrend_df[['supertrend']])
    except Exception as e:
        log.error(f'Exception while trying to calculate the supertrend values: {e}', exc_info=True)
        return df

    return df


def bollinger_bands_indicator_values_df(df):
    if df.empty:
        return df

    try:
        indicator_bb = BollingerBands(close=df['close'])
        df['bb_moving_average'] = indicator_bb.bollinger_mavg()
        df['bb_upper_band'] = indicator_bb.bollinger_hband()
        df['bb_lower_band'] = indicator_bb.bollinger_lband()
    except Exception as e:
        log.error(f'Exception while trying to calculate the BollingerBand Values: {e}', exc_info=True)
        return df

    return df


def average_true_range_indicator_values_df(df):
    if df.empty or df.shape[0] < 14:
        return df

    try:
        atr_indicator = AverageTrueRange(df['high'], df['low'], df['close'])
        df['atr'] = atr_indicator.average_true_range()
    except Exception as e:
        log.error(f'Exception while trying to calculate the AverageTrueRange Values: {e}', exc_info=True)
        return df

    return df


def rsi_indicator_values_df(df, window):
    if df.empty:
        return df

    try:
        rsi_window = RSIIndicator(close=df['close'], window=window)
        df[f'rsi_{window}'] = rsi_window.rsi()
    except Exception as e:
        log.error(f'Exception while trying to calculate the RSI Indicator Values: {e}', exc_info=True)
        return df

    return df


def ema_indicator_values_df(df, window_list):
    if df.empty:
        return df

    try:
        for window in window_list:
            ema_indicator_window = EMAIndicator(close=df['close'], window=window)
            df[f'ema_{window}'] = ema_indicator_window.ema_indicator()
    except Exception as e:
        log.error(f'Exception while trying to calculate the EMAIndicator Values: {e}', exc_info=True)
        return df

    return df


def adx_indicator_values_df(df, window):
    if df.empty or df.shape[0] < 14:
        return df

    try:
        adx = ADXIndicator(df['high'], df['low'], df['close'], window=window)
        df['adx'] = adx.adx()
        df['adx_pos'] = adx.adx_pos()
        df['adx_neg'] = adx.adx_neg()
    except Exception as e:
        log.error(f'Exception while trying to calculate the ADXIndicator Values: {e}', exc_info=True)
        return df

    return df


def sma_indicator_values_df(df, window):
    if df.empty:
        return df

    try:
        sma_indicator_window = SMAIndicator(close=df['close'], window=window)
        df[f'sma_{window}'] = sma_indicator_window.sma_indicator()
    except Exception as e:
        log.error(f'Exception while trying to calculate the SMAIndicator Values: {e}', exc_info=True)
        return df

    return df


def psar_indicator_values_df(df):
    if df.empty:
        return df

    try:
        psar_indicator = PSARIndicator(high=df['high'], low=df['low'], close=df['close'])
        df['psar'] = psar_indicator.psar()
    except Exception as e:
        log.error(f'Exception while trying to calculate the PSARIndicator Values: {e}', exc_info=True)
        return df

    return df


def slope_calc(df, indicator_value, timeperiod):
    if df.empty:
        return df

    try:
        df[f'{indicator_value}_slope'] = pd_ta.slope(df[indicator_value], timeperiod=timeperiod)
    except Exception as e:
        log.error(f'Exception while trying to calculate the Supertrend Buy Sell Signal Values: {e}', exc_info=True)
        return df

    return df


def supertrend_7_3_alpha_signals(df):
    if df.empty:
        return df

    try:
        df['supertrend_7_3_long_entry'] = np.logical_and(df.close.shift(1) <= df.supertrend.shift(1),
                                                         df.close > df.supertrend)
        df['supertrend_7_3_short_entry'] = np.logical_and(df.close.shift(1) >= df.supertrend.shift(1),
                                                          df.close < df.supertrend)
        df['supertrend_7_3_long_exit'] = df['supertrend_7_3_short_entry']
        df['supertrend_7_3_short_exit'] = df['supertrend_7_3_long_entry']
    except Exception as e:
        log.error(f'Exception while trying to calculate the Supertrend Buy Sell Signal Values: {e}', exc_info=True)
        return df

    return df


def adx_ema_alpha_signals(df):
    if df.empty:
        return

    try:
        df['ema_positive_trend'] = np.where(df[FASTER_EMA_WINDOW_COL] > df[SLOWER_EMA_WINDOW_COL], True, False)
        df['adx_positive_trend'] = np.where(df['adx_pos'] > df['adx_neg'], True, False)
        df['close_sma_trend'] = np.where(df['close'] > df[SMA_TREND_INDICATOR_COL], True, False)
        df['close_ema_trend'] = np.where(df['close'] > df[EMA_TREND_INDICATOR_COL], True, False)
        df['rsi_14_slope_trend'] = np.where(df['rsi_14_slope'] > 0, True, False)
        df['adx_pos_slope_trend'] = np.where(df['adx_pos_slope'] > 0, True, False)

        df['adx_ema_strategy_entry'] = np.where(
            (df[FASTER_EMA_WINDOW_COL] > df[SLOWER_EMA_WINDOW_COL]) & (df['adx_pos'] > df['adx_neg']), True, False)

        df['adx_ema_strategy_verify'] = np.where((df['close'] > df[SMA_TREND_INDICATOR_COL]) &
                                                 (df['close'] > df[EMA_TREND_INDICATOR_COL]) &
                                                 (df['rsi_14_slope'] > 0) &
                                                 (df['adx_pos_slope'] > 0), True, False)

        df['adx_ema_long_entry'] = False
        adx_ema_strategy_entry_group = False
        for index, row in df.iterrows():
            if not row['adx_ema_strategy_entry']:
                adx_ema_strategy_entry_group = False

            if row['adx_ema_strategy_entry'] == 1 and row[
                'adx_ema_strategy_verify'] and not adx_ema_strategy_entry_group:
                df.at[index, 'adx_ema_long_entry'] = True
                adx_ema_strategy_entry_group = True

        df['adx_ema_long_exit'] = False
        df['adx_ema_short_entry'] = False
        df['adx_ema_short_exit'] = False
    except Exception as e:
        log.error(f'Exception while to calculate the ADX EMA Strategy Alpha Signals: {e}', exc_info=True)
        return df

    return df


def adx_ema_basic_alpha_signals(df):
    if df.empty:
        return df

    try:
        df['adx_ema_basic_ema_dmi_positive_trend'] = np.where(
            (df[FASTER_EMA_WINDOW_COL] > df[SLOWER_EMA_WINDOW_COL]) & (df['adx_pos'] > df['adx_neg']), 1, 0)
        df['adx_ema_basic_ema_dmi_positive_trend_diff'] = df['adx_ema_basic_ema_dmi_positive_trend'].diff().fillna(
            0).astype(int)
        df['adx_ema_basic_long_entry'] = np.where((df['adx_ema_basic_ema_dmi_positive_trend_diff'] == 1), True, False)
        df['adx_ema_basic_long_exit'] = False

        df['adx_ema_basic_ema_dmi_negative_trend'] = np.where(
            (df[FASTER_EMA_WINDOW_COL] < df[SLOWER_EMA_WINDOW_COL]) & (df['adx_pos'] < df['adx_neg']), 1, 0)
        df['adx_ema_basic_ema_dmi_negative_trend_diff'] = df['adx_ema_basic_ema_dmi_negative_trend'].diff().fillna(
            0).astype(int)
        df['adx_ema_basic_short_entry'] = np.where((df['adx_ema_basic_ema_dmi_negative_trend_diff'] == 1), True, False)
        df['adx_ema_basic_short_exit'] = False
    except Exception as e:
        log.error(f'Exception while trying to generate the EMA Basic Long Short Strategy Signal: {e}', exc_info=True)
        return df

    return df


def ichimoku_chikou_span_alpha_signals_basic(df):
    if df.empty:
        return df

    try:
        df['ichimoku_chikou_span_basic_long_entry'] = np.logical_and((df.ICS_26.shift(1) < df.low.shift(1)),
                                                                     df.ICS_26 > df.high)

        df['ichimoku_chikou_span_basic_short_entry'] = np.logical_and((df.ICS_26.shift(1) > df.high.shift(1)),
                                                                      df.ICS_26 < df.low)

        df['ichimoku_chikou_span_basic_long_exit'] = df['ichimoku_chikou_span_basic_short_entry']
        df['ichimoku_chikou_span_basic_short_exit'] = df['ichimoku_chikou_span_basic_long_entry']
    except Exception as e:
        log.error(f'Exception while trying to Ichimoku Chikou Span Basic Buy/Sell Signal: {e}', exc_info=True)
        return df

    return df


def ichimoku_chikou_span_alpha_signals_20_days_high_low_crossover(df):
    if df.empty:
        return df

    try:
        df['ichimoku_chikou_span_20_days_high_low_crossover_long_entry'] = np.logical_and(
            (df.ICS_26.shift(1) < df.low.shift(1)) &
            (df.ICS_26.shift(2) < df.low.shift(2)) &
            (df.ICS_26.shift(3) < df.low.shift(3)) &
            (df.ICS_26.shift(4) < df.low.shift(4)) &
            (df.ICS_26.shift(5) < df.low.shift(5)) &
            (df.ICS_26.shift(6) < df.low.shift(6)) &
            (df.ICS_26.shift(7) < df.low.shift(7)) &
            (df.ICS_26.shift(8) < df.low.shift(8)) &
            (df.ICS_26.shift(9) < df.low.shift(9)) &
            (df.ICS_26.shift(10) < df.low.shift(10)) &
            (df.ICS_26.shift(11) < df.low.shift(11)) &
            (df.ICS_26.shift(12) < df.low.shift(12)) &
            (df.ICS_26.shift(13) < df.low.shift(13)) &
            (df.ICS_26.shift(14) < df.low.shift(14)) &
            (df.ICS_26.shift(15) < df.low.shift(15)) &
            (df.ICS_26.shift(16) < df.low.shift(16)) &
            (df.ICS_26.shift(17) < df.low.shift(17)) &
            (df.ICS_26.shift(18) < df.low.shift(18)) &
            (df.ICS_26.shift(19) < df.low.shift(19)),
            df.ICS_26 > df.high)

        df['ichimoku_chikou_span_20_days_high_low_crossover_short_entry'] = np.logical_and(
            (df.ICS_26.shift(1) > df.high.shift(1)) &
            (df.ICS_26.shift(2) > df.high.shift(2)) &
            (df.ICS_26.shift(3) > df.high.shift(3)) &
            (df.ICS_26.shift(4) > df.high.shift(4)) &
            (df.ICS_26.shift(5) > df.high.shift(5)) &
            (df.ICS_26.shift(6) > df.high.shift(6)) &
            (df.ICS_26.shift(7) > df.high.shift(7)) &
            (df.ICS_26.shift(8) > df.high.shift(8)) &
            (df.ICS_26.shift(9) > df.high.shift(9)) &
            (df.ICS_26.shift(10) > df.high.shift(10)) &
            (df.ICS_26.shift(11) > df.high.shift(11)) &
            (df.ICS_26.shift(12) > df.high.shift(12)) &
            (df.ICS_26.shift(13) > df.high.shift(13)) &
            (df.ICS_26.shift(14) > df.high.shift(14)) &
            (df.ICS_26.shift(15) > df.high.shift(15)) &
            (df.ICS_26.shift(16) > df.high.shift(16)) &
            (df.ICS_26.shift(17) > df.high.shift(17)) &
            (df.ICS_26.shift(18) > df.high.shift(18)) &
            (df.ICS_26.shift(19) > df.high.shift(19)),
            df.ICS_26 < df.low)

        df['ichimoku_chikou_span_20_days_high_low_crossover_long_exit'] = df['ichimoku_chikou_span_20_days_high_low_crossover_short_entry']
        df['ichimoku_chikou_span_20_days_high_low_crossover_short_exit'] = df['ichimoku_chikou_span_20_days_high_low_crossover_long_entry']
    except Exception as e:
        log.error(f'Exception while trying to Ichimoku Chikou Span 20 Days High Low Crossover Buy/Sell Signal: {e}', exc_info=True)
        return df

    return df


def ichimoku_chikou_span_alpha_signals(df):
    if df.empty:
        return df

    try:
        df['ichimoku_chikou_span_long_entry'] = np.logical_and((df.ICS_26.shift(1) <= df.close.shift(1)) &
                                                               (df.ICS_26.shift(2) < df.close.shift(2)) &
                                                               (df.ICS_26.shift(3) < df.close.shift(3)) &
                                                               (df.ICS_26.shift(4) < df.close.shift(4)) &
                                                               (df.ICS_26.shift(5) < df.close.shift(5)) &
                                                               (df.ICS_26.shift(6) < df.close.shift(6)),
                                                               df.ICS_26 > df.close)

        df['ichimoku_chikou_span_short_entry'] = np.logical_and((df.ICS_26.shift(1) >= df.close.shift(1)) &
                                                                (df.ICS_26.shift(2) > df.close.shift(2)) &
                                                                (df.ICS_26.shift(3) > df.close.shift(3)) &
                                                                (df.ICS_26.shift(4) > df.close.shift(4)) &
                                                                (df.ICS_26.shift(5) > df.close.shift(5)) &
                                                                (df.ICS_26.shift(6) > df.close.shift(6)),
                                                                df.ICS_26 < df.close)

        df['ichimoku_chikou_span_long_exit'] = df['ichimoku_chikou_span_short_entry']
        df['ichimoku_chikou_span_short_exit'] = df['ichimoku_chikou_span_long_entry']
    except Exception as e:
        log.error(f'Exception while trying to Ichimoku Chikou Span Buy/Sell Signal: {e}', exc_info=True)
        return df

    return df


def ichimoku_chikou_span_strong_alpha_signals(df):
    if df.empty:
        return df

    try:
        df['ichimoku_chikou_span_strong_long_entry'] = np.logical_and(
            (df.ICS_26.shift(1) <= df.close.shift(1)) &
            (df.ICS_26.shift(2) < df.close.shift(2)) &
            (df.ICS_26.shift(3) < df.close.shift(3)) &
            (df.ICS_26.shift(4) < df.close.shift(4)) &
            (df.ICS_26.shift(5) < df.close.shift(5)) &
            (df.ICS_26.shift(6) < df.close.shift(6)),
            df.ICS_26 > df.high)

        df['ichimoku_chikou_span_strong_short_entry'] = np.logical_and(
            (df.ICS_26.shift(1) >= df.close.shift(1)) &
            (df.ICS_26.shift(2) > df.close.shift(2)) &
            (df.ICS_26.shift(3) > df.close.shift(3)) &
            (df.ICS_26.shift(4) > df.close.shift(4)) &
            (df.ICS_26.shift(5) > df.close.shift(5)) &
            (df.ICS_26.shift(6) > df.close.shift(6)),
            df.ICS_26 < df.low)

        df['ichimoku_chikou_span_strong_long_exit'] = df['ichimoku_chikou_span_strong_short_entry']
        df['ichimoku_chikou_span_strong_short_exit'] = df['ichimoku_chikou_span_strong_long_entry']
    except Exception as e:
        log.error(f'Exception while trying to Ichimoku Chikou Span Buy/Sell Signal Strong: {e}', exc_info=True)
        return df

    return df


def ichimoku_chikou_span_strongest_alpha_signals(df):
    if df.empty:
        return df

    try:
        df['ichimoku_chikou_span_strongest_long_entry'] = np.logical_and(
            (df.ICS_26.shift(1) <= df.low.shift(1)) &
            (df.ICS_26.shift(2) < df.low.shift(2)) &
            (df.ICS_26.shift(3) < df.low.shift(3)) &
            (df.ICS_26.shift(4) < df.low.shift(4)) &
            (df.ICS_26.shift(5) < df.low.shift(5)) &
            (df.ICS_26.shift(6) < df.low.shift(6)) &
            (df.ICS_26.shift(7) < df.low.shift(7)) &
            (df.ICS_26.shift(8) < df.low.shift(8)),
            df.ICS_26 > df.high)

        df['ichimoku_chikou_span_strongest_short_entry'] = np.logical_and(
            (df.ICS_26.shift(1) >= df.high.shift(1)) &
            (df.ICS_26.shift(2) > df.high.shift(2)) &
            (df.ICS_26.shift(3) > df.high.shift(3)) &
            (df.ICS_26.shift(4) > df.high.shift(4)) &
            (df.ICS_26.shift(5) > df.high.shift(5)) &
            (df.ICS_26.shift(6) > df.high.shift(6)) &
            (df.ICS_26.shift(7) > df.high.shift(7)) &
            (df.ICS_26.shift(8) > df.high.shift(8)),
            df.ICS_26 < df.low)

        df['ichimoku_chikou_span_strongest_long_exit'] = df['ichimoku_chikou_span_strongest_short_entry']
        df['ichimoku_chikou_span_strongest_short_exit'] = df['ichimoku_chikou_span_strongest_long_entry']
    except Exception as e:
        log.error(f'Exception while trying to Ichimoku Chikou Span Buy/Sell Signal Strongest: {e}', exc_info=True)
        return df

    return df


if __name__ == '__main__':
    start_time = time.time()
    parser = argparse.ArgumentParser()
    parser.add_argument('limitRowsFlag')
    args = parser.parse_args()
    access_token = AliceBlue.login_and_get_access_token(username=alice_blue_username, password=alice_blue_password,
                                                        twoFA=alice_blue_twoFA, api_secret=alice_blue_api_secret)
    alice_blue = AliceBlue(username=alice_blue_username, password=alice_blue_password, access_token=access_token,
                           master_contracts_to_download=['NSE', 'BSE', 'NFO'])

    alice_blue_client = AliceBlueDailySignals(alice_blue, (args.limitRowsFlag == 'True'))
    alice_blue_client.start_alice_blue_daily_signals_generation()
    total_time_taken = time.time() - start_time
    log.info("=====================================================================")
    log.info("Finished Generation of Alice Blue Daily Signals!")
    log.info("Time Taken: --- %s seconds --- (%s minutes)" % (total_time_taken, total_time_taken/60))
    log.info("=====================================================================")