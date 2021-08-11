def get_value(nestedObj, key):
    key_path = key.split("/")
    res = nestedObj
    for item in key_path:
        res = res[item]
    return res


nestedObject = {"a":{"b":{"c":"d"}}}
key = "a/b/c"
print(get_value(nestedObject, key))

nestedObject = {"x":{"y":{"z":"a"}}}
key = "x/y/z"
print(get_value(nestedObject, key))
