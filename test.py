import requests

def test1():
    
    r = requests.put(url="http://127.0.0.1:8080/google/?v=monkey")
    print("server answer = ", r.status_code)
    
    r = requests.get(url="http://127.0.0.1:8080/google")
    print(r.content)
    

def test2():
    r = requests.put(url="http://127.0.0.1:8080/teacher/?v=Look")
    print("server answer = ", r.status_code)
    
    
def test3():
    r = requests.get(url="http://127.0.0.1:8080/time/?timeout=5")
    print(r.status_code)
    
def test4():
    r = requests.put(url="http://127.0.0.1:8080/time/?v=monkey")

def main():
    test1()
    test2()
    test3()
    test4()

if __name__ == "__main__":
    main()