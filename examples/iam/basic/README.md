# Basic IAM User

This example provides the following resources:

* Create an identity user, name as 'user_A'.
* Two ways to create identity group:
    + Create a single identity group(defalut group).
    + Create multiple identity groups from a list(Object({string,string})).
* Add user_A to second identity group. Add user_A to the default identity group when second custom identity groups do
  not exist.
* Assign permissions to roles in the identity group as an administrator(By domain).
