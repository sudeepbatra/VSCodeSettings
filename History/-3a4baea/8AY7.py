from py5paisa import FivePaisaClient

client = FivePaisaClient(email="random_email@xyz.com", passwd="password", dob="YYYYMMDD")
client.login()
```

After successful authentication, you should get a `Logged in!!` message.