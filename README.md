# Camel
Camel is an interpreted programming language written entirely in Go without use of any third-party libraries. It uses pratt parsing to create AST (abstract syntax tree). Evaluation is done simply by walking the tree. It is supposed to be simple and easy to understand.  Read Acknowledgements for more information.
## Installation 
Download `repl` at /bin to run the interpreter. 
```./repl


     __,  .-.  .-.
    (__ _ |  \/  | _
       (_|||\__/||(/_|
          ||    ||   |_
         _||   _||
         ""    ""

------------------------------------------------
Hello! mehrdad, This is camel programming language!
Feel free to type in commands
>> 
```
### Variables
Camel is dynamically-typed. Decalare variables using `beza` keyword.
```rust 
>> beza x = 2 
>> beza y = "name" 
>> beza t = true 
>> x + 3 
5
>> !t 
false 
>> t == (1 < 2) 
true
>> "name" + " : " + "Monica" 
name : Monica 
```
### Array 
```rust 
>> beza x = [1 , 2 , "hey", true]
>> x[5 - 3*2 + 1] + x[1] 
3
```
### Hash
```rust 
>> beza x = {"lang": "camel" , "version": 0.0}
>> x["lang"] 
camel
```
### Builtins 
```rust 
>> beza x = [1 , 2 , "hey", true]
>> len(x)
4 
>> peek(x) 
true 
>> pop(x) 
[1, 2, "hey"] 
```
### Condition
```rust
if ( 2 - 4 < 0 ) {
  2 
} else { 
  bede 1 
}
```
Note that you can choose to write `bede` or skip. same goes with semicolons. `bede` is a keyword used to return values.
### Function 
keyword for functions are `foo` & `bar`. use any of these two to define a function.
```rust 
beza fib = foo(x) { 
  if (x == 0) {
    0
  } else {
    if (x == 1) {
      bede 1;
    } else {
      fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};
```
### Closure
```rust 
>> beza newAddr = foo(x) { bar(c) { x + c } }
>> beza addTwo  = newAddr(2)
>> addTwo(5) 
7
>> newAddr(2)(5) 
7
```
### Errors
```rust
>> beza x = 2 
>> x()
Error: Invalid function call, INTEGER is not a function
>> beza x = [1 , 2] 
>> x[4]
Error: Index out of range
>> x["hey"] 
Error: Invalid Index: index operator not supported for type ARRAY
>> 2 == true
Error: Type mismatch: invalid operator == for types INTEGER BOOLEAN
```

## TODO 

- [ ] Add support for bitwise operators 
- [ ] Add support for logical operators 
- [ ] Add support for modulo operators 
- [ ] Add support for emojis 
- [ ] Resolve hash collisions
- [ ] Scanning input
- [ ] Add support for comments
- [ ] Add support for control characters
- [ ] Add support for loops
- [ ] Add support for error handling in parser 



## Acknowledgement 
The Camel programming language is highly based on Monkey. Moneky is an educational programming language first presented in the book "Writing an interpreter in Go". Camel was built following instructions of the book. You can find more information [here](https://interpreterbook.com/).
