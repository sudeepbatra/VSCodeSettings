from breeze_connect import BreezeConnect

app_name = 'StoxMarket'
api_key = 's76162#+U35414Y*S413=099_FA6P567'
secret_key = 'I04M0Y9!5vP7G3ct53e395+41F27=621'

breeze = BreezeConnect(api_key=api_key)

import urllib
print("https://api.icicidirect.com/apiuser/login?api_key="+urllib.parse.quote_plus(api_key))

#From the stox.market network request payload pick the api_session value:
api_session = 16963048

breeze.generate_session(api_secret=secret_key, session_token=api_session)
 
customer_details = breeze.get_customer_details(api_session=api_session)
print(customer_details)
# Prints
#{'Success': {'exg_trade_date': {'NSE': '22-Jul-2022', 'BSE': '22-Jul-2022', 'FNO': '22-Jul-2022', 'NDX': '22-Jul-2022'}, 'exg_status': {'NSE': 'C', 'BSE': 'C', 'FNO': 'C', 'NDX': 'C'}, 'segments_allowed': {'Trading': 'Y', 'Equity': 'Y', 'Derivatives': 'Y', 'Currency': 'N'}, 'idirect_userid': 'SUDDEPBA', 'idirect_user_name': 'SUDEEP   BATRA', 'idirect_ORD_TYP': '', 'idirect_lastlogin_time': '22-Jul-2022 08:18:44', 'mf_holding_mode_popup_flg': 'Y', 'commodity_exchange_status': 'C', 'commodity_trade_date': '22-Jul-2022', 'commodity_allowed': 'C'}, 'Status': 200, 'Error': None}

funds = breeze.get_funds()
print(funds)
# Prints
# {'Success': {'bank_account': '000201066633', 'total_bank_balance': 112.29, 'allocated_equity': 112.29, 'allocated_fno': 0.0, 'block_by_trade_equity': 0.0, 'block_by_trade_fno': 0.0, 'block_by_trade_balance': 0.0, 'unallocated_balance': '2.16'}, 'Status': 200, 'Error': None}
