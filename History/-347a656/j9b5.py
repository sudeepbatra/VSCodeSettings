from smartapi import SmartConnect

#create object of call
obj=SmartConnect(api_key="zlndzd6z")

data = obj.generateSession("S1632585","DevDas123")
refreshToken= data['data']['refreshToken']

feedToken=obj.getfeedToken()

userProfile = obj.getProfile(refreshToken)
print(userProfile)

try:
    historicParam={
    "exchange": "NSE",
    "symboltoken": "3045",
    "interval": "ONE_MINUTE",
    "fromdate": "2021-02-08 09:00",
    "todate": "2021-02-08 09:16"
    }
    candle_data = obj.getCandleData(historicParam)
    print(candle_data)
except Exception as e:
    print("Historic Api failed: {}".format(e.message))
# #logout
# try:
#     logout=obj.terminateSession('Your Client Id')
#     print("Logout Successfull")
# except Exception as e:
#     print("Logout failed: {}".format(e.message))

#place order
try:
    orderparams = {
        "variety": "NORMAL",
        "tradingsymbol": "SBIN-EQ",
        "symboltoken": "3045",
        "transactiontype": "BUY",
        "exchange": "NSE",
        "ordertype": "LIMIT",
        "producttype": "INTRADAY",
        "duration": "DAY",
        "price": "19500",
        "squareoff": "0",
        "stoploss": "0",
        "quantity": "1"
        }
    orderId=obj.placeOrder(orderparams)
    print("The order id is: {}".format(orderId))
except Exception as e:
    print("Order placement failed: {}".format(e.message))