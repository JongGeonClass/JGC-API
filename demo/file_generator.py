import os
import glob

from shutil import copyfile

def main():
    print('delete all description files')
    [os.remove(f) for f in glob.glob("./product/description/*")]
    print('delete all title files')
    [os.remove(f) for f in glob.glob("./product/title/*")]

    os.mkdir('./product')
    os.mkdir('./product/description')
    print('generate description files')
    for i in range(1, 51):
        with open(f'./product/description/{i}.txt', 'w') as f:
            f.write(f'''# 종건급 Product {i} description

진자 개지리는 종건급 상품입니다.

와우 너무 개지림

와 진짜 이거 왜 안씀??

님들 진자 후회하는거임

종건급 상품 개지립니다.

우주 최강 종건급 상품 ~~''')
    
    os.mkdir('./product/title')
    print('generate title files')
    for i in range(1, 51):
        copyfile('./testimg.png', f'./product/title/{i}.png')

main()