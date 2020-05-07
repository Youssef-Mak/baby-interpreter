# Baby is a small, functional-based programming language

[![Run on Repl.it](https://repl.it/badge/github/Youssef-Mak/baby-interpreter)](https://repl.it/@YoussefMak1/baby-interpreter)

Think of simplified version of JavaScript with a more welcoming syntax all while limiting ambiguity and providing powerful built-in features. The interpreter is heavily influenced by the one
built in Thorsten Ball's books(Monkey) but it aims to extend the language by providing various enhancements such as different levels of assignment(deep-copy, shallow-copy etc), operators, and a wider array of built-in functionality. (while loops, doWhile, etc)

Below is a functional implementation of MergeSort in Baby.

<details>
<summary>Functional MergeSort in Baby</summary>
<p>

```ocaml
let funcMS = fun(cmp, x, y) {
	if (isEmpty(x) & (!isEmpty(y))) {
		return y;
	}
	if ((!isEmpty(x)) & isEmpty(y)) {
		return x;
	}

	let result = [];
	if (cmp(head(x), head(y))) {
		result = append(result, head(x));
		result = append(result, funcMS(cmp, rest(x), y));
	} else {
		result = append(result, head(y));
		result = append(result, funcMS(cmp, x, rest(y)));
	};
	return result;
};

let split = fun(x, y, z) {
	let parts = [0,1];
	if (isEmpty(x)) {
		parts = insert(parts, y, 0);
		parts = insert(parts, z, 1);
		return parts;
	} else {
		return split(rest(x), z, append(head(x),y));
	};
};

let mergesort = fun(cmp, x) {
	if (isEmpty(x) | isEmpty(rest(x))) {
		return x;
	} else {
		let parts = split(x, [], []);
		return funcMS(cmp, mergesort(cmp, get(parts, 0)), mergesort(cmp, get(parts, 1)));
	};
};

mergesort(fun(x,y) { return (x<y);}, [12, 11, 13, 5, 6, 7]) =*= [5,6,7,11,12,13];
```

</p>
</details>

- **Baby is small.** The Interpreter implementation is around [3,500 lines][src].
  Making it understandable and maintainable.

- **Baby is reasonably fast.** For an interpreted language, Baby is reasonably fast. With GO being the host language, it makes use of GoLang's performant Garbage Collector
  and speedy memory allocation.

* **Baby encourages declarative programming.** Imperative programming can be confusing to beginners.
  Baby aims to clear the confusion by providing a declarative safe space while providing imperative tools
  so everyone is happy. :)

If you like the sound of this? You can even try
it [in your browser][browser]!

# Language Tour

Baby's syntax is designed to be familiar to people coming from JavaScript-like languages
while being a lot simpler.

Code is stored in plain text files with a `.bb` file extension. Baby is an interpreted language and is
not compiled

## Reserved words

Here are the reserved keywords in Baby:

    fun
    let
    if
    else
    while
    return

## Operators

Other than the expected operators that are included in most programming languages, Baby throws some new operators to the mix:

### Special Equality Operators

#### Reference Equality

Evaluating whether references are identical is denoted by the infix operator `=&=`.

Likewise, reference inequality is denoted by `!&=`.

#### Value Equality

Evaluating whether values at addresses are identical is denoted by the infix operator `=*=`.

Likewise, value inequality is denoted by `!*=`.

### Special Assignment Operators

#### Reference Assignment

In Baby, making a deep copy is simple. In this case, after assignment both the value and the address are assigned.
In Baby this is denoted by the `=&` infix operator.

      b =& a // Shallow copy (assign the value and the reference)
      b =&= a ----> true
      b =\*= a ----> true

#### Value Assignment

In Baby, making a deep copy is simple. In this case, deep copy entails that after assignment only the value is assigned
not the address. In Baby this is denoted by the `=*` infix operator and can be inferred with the `=` infix operator.

Here is an example:

      b =* a // Deep copy (simply take the value)
      b =&= a ----> false
      b =*= a ----> true

## Identifiers

Identifiers must be composed of letters and can contain underscores. CamelCase or kebab-case are encouraged.
Case is sensitive. Numbers in identifiers(ex: foo3) are not supported.

    hi
    camelCase
    PascalCase
    _under_score
    ALL_CAPS

## Blocks

### While Loops

While loops are implements just as most other programming languages:

`while (<condition>) {<consequence>}`

### Do while Loops

Do while loops are a little trickier: first, a function that returns a boolean entailing the condition of the loop
should be defined. The body of this function will be the body of the loop. Then a call is made to the built-in
doWhile function which does all the magic.

```
let i = 6;
let inc = fun() {i=i+1; return (i<5)};
doWhile(inc);
print(i); // 7
```

### If-Else statements

If-Else statement are structured like so: `if (<condition>) <consequence> else <alternative>`

### Functions

Functions are declared like so : `let <identifier> = fun(<list of params(identifiers) {<list of statements>}`

Baby supports closures as well as the passing of functions(higher-order functions).
The return keyword can be omitted but is recommended for code readability.

## Declaration Statements

Initial Assignment is done with `let` like so: `let <identifier> = <expression>`.
In re-assignment `let` can be omitted like so `<identifier> = <expression>`.

[browser]: https://repl.it/@YoussefMak1/baby-interpreter
[src]: https://github.com/Youssef-Mak/baby-interpreter/tree/master/pkg
