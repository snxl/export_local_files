tempPath = "./tmp"

for x in range(1,5000):
   f = open(tempPath + "/test"+str(x)+".txt", 'w')
   f.write(str(x) + "\n" + str(5000-x))
   f.close()