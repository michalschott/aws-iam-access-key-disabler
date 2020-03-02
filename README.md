# aws-iam-access-key-disabler

## why?

Usually your engineers do not care about best security practices due to enormous amount of work they need to deliver. I want to make your life easier, and at the same time upskill myself in GO and Serverless.

## standalone tool

Just standalone tool you can run from anywhere.

## aws lambda terraform v11 module

Haven't figured out how to ship go lambdas properly, so temporary workaround:

```
cd aws-lambda
make
```

Now you can include terraform module to your base.

## ROADMAP
- add travis / circleci integrations
- update docs
- ...
