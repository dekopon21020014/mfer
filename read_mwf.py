f = open('sample-data/ECG01.mwf', 'rb')
data = f.read()
print(data)
for i in range(256):
    print(hex(data[i]))
