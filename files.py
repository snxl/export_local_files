import os


tempPath = "./tmp"
filesToGenerate = 5000

if not os.path.exists(tempPath):
   os.makedirs(tempPath)

for x in range(1,filesToGenerate):
   f = open(tempPath + "/test"+str(x)+".txt", 'w')
   f.write(str(x) + "\n" + str(filesToGenerate-x))
   f.close()