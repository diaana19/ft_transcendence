# Code Commenting Guidelines

To maintain clean and readable code, comments must follow strict placement and content rules.

---

## General Rules

* Comments must be written in English
* Comments must be clear, concise, and relevant
* Do not write comments for obvious code
* Avoid redundant or trivial comments

---

## Allowed Comment Locations

### 1. File Header

Each file must begin with a header comment describing its purpose.

Content must include:

* File purpose
* Main responsibilities
* Any important notes (if necessary)

Example:

```
/*
** File: parser.c
** Description: Handles input parsing and validation
** Responsibilities:
** - Parse raw input
** - Validate syntax
** - Return structured data
*/
```

---

### 2. Function-Level Comments

Each function must have a comment placed directly above it.

Content must include:

* Purpose of the function
* Parameters
* Return value

Example:

```
/*
** Parses user input and returns a structured object
** params:
**   input (char*): raw user input
** returns:
**   parsed structure or NULL on failure
*/
```

---

## Forbidden Practices

* Do not place comments inside function bodies
* Do not comment individual lines of code
* Do not explain what the code is doing step-by-step
* Do not leave commented-out code

Bad example:

```
x++; // increment x
```

---

## Exceptions

Inline comments are allowed only if:

* The logic is complex and cannot be easily understood
* A non-obvious decision must be explained

These cases must remain rare.

Example:

```
/* Intentional offset to match external protocol specification */
index += 4;
```

---

## Notes

* Code should be self-explanatory whenever possible
* Prefer meaningful variable and function names over comments
* Poorly written comments may be removed during review
