from py5paisa import FivePaisaClient

# This is not working
cred = {
    "APP_NAME": "5P50603710",
    "APP_SOURCE": "10427",
    "USER_ID": "uxuZEFys5nv",
    "PASSWORD": "7elTHyW0EC3",
    "USER_KEY": "sR12m8nkT8VEPXtfgLFlspj5BQlSqB51",
    "ENCRYPTION_KEY": "jTS6yEtvhXThvDTYNHQNVXmklWFEaeQj"
}

client = FivePaisaClient(email="sudeep.batra@gmail.com", passwd="Hagar123", dob="19771013", cred=cred)
client.login()
print(client)

a = {
    "method": "subscribe",
    "operation": "20depth",
    "instruments": ["NC2885"]
}
print(client.socket_20_depth(a))


def on_message(ws, message):
    print(message)


client.receive_data(on_message)
