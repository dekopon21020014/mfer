f = open('ECG01.mwf', 'rb')
data = f.read()
hoge = data[1]
print(hoge)
