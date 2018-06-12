import fileinput

def hex(c):
    return tuple(int(c[1 + i * 2 : 1 + i * 2 + 2], 16) for i in range(3))

with fileinput.input(files="temp", inplace=True, backup=".bak") as f:
    for line in f:
        temp = str(line).split()
        color = hex(temp[1])
        print("\"" + temp[0] + "\" : color.RGBA{" + str(color[0]) + ", " + str(color[1])  + ", " + str(color[2]) + ", 255}, ")