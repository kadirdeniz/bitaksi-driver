# bitaksi-driver
This is a driver service repository for bitaksi case study.
You can find the `Matcher Service` respository link in below

https://github.com/kadirdeniz/bitaksi-matcher

[![architecture](https://www.linkpicture.com/q/Screen-Shot-2023-01-13-at-11.19.15.png)](https://www.linkpicture.com/view.php?img=LPic63c11b745529f1064691485)

## Table Of Contents
* [General info](#general-info)
* [Clone the project](#clone-the-project)
* [Setup](#setup)
* [Test](#test)

## General info
This service's responsibility is handling `mongodb` operations, when the application starts, firstly it will seed mongodb with given `Coordinated.csv` then can be used for `finding nearest driver` and `batch operations`.This service can only used by `Matcher Service` so i created a `api_key` field for every request, if the api_key is wrong or not provided service will be response with `error:invalid api key`.In this project i tried to follow `Test Driven Development`.`TDD` is a software development approach which test should write before code,then refactoring should be done. I implemented `pactflow` contract test, integration and unit tests and i used `Ginkgo` for clear understanding of test cases. There is a `postman` collection under the `/docs` folder, it can be used for testing.

## Clone the project
```
$ git clone https://github.com/kadirdeniz/bitaksi-driver.git
$ cd bitaksi-driver
```

## Setup
##### Application can run with using docker but running locally would be more stable

```
$ make run

$ make dockerized
$ make run-docker
```


 ## Test
 ```
 $ make tests
 ```
