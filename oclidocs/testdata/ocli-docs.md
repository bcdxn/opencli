<!-- Markdown generated by ocli-codegen DO NOT EDIT. -->

# Pleasantries
![ocli-badge](https://img.shields.io/badge/OpenCLI_Spec-Compliant-brightgreen?link=https%3A%2F%2Fgithub.com%2Fbcdxn%2Fopencli)

_A fun CLI to greet or bid farewell_

---

## Commands


### `$ pleasantries {command} <name> [flags]`

`group` `nonexecutable`


---


### `$ pleasantries farewell <name> [flags]`

Say goodbye
#### Arguments

##### `<name>`

A name to include in the farewell

`string`

#### Flags

##### `--language`

The language of the greeting

`string` `enum` `default:english`

###### Supported Values of `--language`

- `english`
- `spanish`

---


### `$ pleasantries greet <name> [flags]`

Say hello
#### Arguments

##### `<name>`

A name to include the greeting

`string`

#### Flags

##### `--language`

The language of the greeting

`string` `enum` `default:english`

###### Supported Values of `--language`

- `english`
- `spanish`

---

<div style="text-align:center;font-size:12px;">generated by <a href="https://github.com/bcdxn/opencli">OpenCLI</a></div>


