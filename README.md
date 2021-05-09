![](https://github.com/grines/scour/blob/main/scour.gif)

# Scour - AWS Exploitation Framework

Attack AWS

## Status
- This is still an active work in progress. **Lost of Bugs
- Not ready for production

- **Where to file issues**:
[https://github.com/grines/scour/issues](https://github.com/grines/Scour/issues)

- **Maintained by**:
[grines](https://github.com/)

# What is Scour?

Scour is a modern module based AWS exploitation framework written in golang, designed for red team testing and blue team analysis. Scour contains modern techniques that can be used to attack environemnts or build detections from attackers.

# Features
- [X] Command Completion
- [X] Dynamic listing resources
- [X] Command history
- [X] Blue team mode (tags attacks with unique User Agent)

## Installation

Scour is written in golang so its easy to ship around as a binary.

##Gettable
go get github.com/grines/scour

##Build
go build main.go


For a more detailed and user-friendly set of user instructions, please check out the Wiki's [installation guide](https://github.com/grines/scour/wiki/Installation).

## Scour's Modules

Scour uses a range modules:
- [X] Operations (2)  
- [X] Enumeration (7)
- [X] Privilege Escalation (3)
- [X] Lateral Movement (2)
- [X] Evasion (5)
- [X] Credential Discovery (4)
- [X] Execution (2)
- [X] Persistance (7)
- [X] Exfiltration (1)

## Notes

* Scour is supported on all Linux/OSX.
* Scour is Open-Source Software and is distributed with a BSD-3-Clause License.

## Getting Started

The first time Scour is launched, 

## Basic Commands in Scour

* `token profile <profile_name> <region>` will list the available aws profiels stored in ~/aws/credentials.
* `help module` will return the applicable help information for the specified module.
* `attack evasion <tactic>` will run the specified module with its default parameters.

## Running Scour From the command line

* `scour` will enter cli mode
* `Not Connected <> token profile apiuser us-east-1` sets the session to use for commands that require one
* `Connected <apiuser/us-east-1>` actively connected to an aws profile from (~,/aws/credentials) in (region)
* `Connected <apiuser/us-east-1> attack enum <attack>` tab completion with list available enumeration tactics
* `Connected <apiuser/us-east-1> attack privesc <attack>` tab completion with list available enumeration tactics
* `Connected <apiuser/us-east-1> attack lateral <attack>` tab completion with list available enumeration tactics
* `Connected <apiuser/us-east-1> attack evasion <attack>` tab completion with list available enumeration tactics
* `Connected <apiuser/us-east-1> attack creds <attack>` tab completion with list available enumeration tactics
* `Connected <apiuser/us-east-1> attack execute <attack>` tab completion with list available enumeration tactics
* `Connected <apiuser/us-east-1> attack persist <attack>` tab completion with list available enumeration tactics
* `Connected <apiuser/us-east-1> attack exfile <attack>` tab completion with list available enumeration tactics

## Enumeration

* `Connected <apiuser/us-east-1> attack enum IAM` IAM discovery
```UA Tracking: exec-env/FhFIm7mvmp/nnK7NmJXNF/iam-enum
+-------------+---------------------+------------------+---------------+--------------+
|    USER     |  MANAGED POLICIES   | INLINE POLICIES  |    GROUPS     | ISPRIVILEGED |
+-------------+---------------------+------------------+---------------+--------------+
| admin  | AdministratorAccess | AllEKSInlineuser | SecurityAudit | true         |
| EC2    | AmazonEC2FullAccess |                  |               | true         |
+-------------+---------------------+------------------+---------------+--------------+```
* `Connected <apiuser/us-east-1> attack enum Roles` Roles discovery
* `Connected <apiuser/us-east-1> attack enum EC2` EC2 discovery
* `Connected <apiuser/us-east-1> attack enum S3` S3 discovery
* `Connected <apiuser/us-east-1> attack enum Groups` Groups discovery
* `Connected <apiuser/us-east-1> attack enum Network` Network discovery

## Disclaimers, and the AWS Acceptable Use Policy

* To the best of our knowledge Scour's capabilities are compliant with the AWS Acceptable Use Policy, but as a flexible and modular tool, we cannot guarantee this will be true in every situation. It is entirely your responsibility to ensure that how you use Scour is compliant with the AWS Acceptable Use Policy.
* Depending on what AWS services you use and what your planned testing entails, you may need to [request authorization from Amazon](https://aws.amazon.com/security/penetration-testing/) before actually running Scour against your infrastructure. Determining whether or not such authorization is necessary is your responsibility.
* As with any penetration testing tool, it is your responsibility to get proper authorization before using Scour outside of your environment.
* Scour is software that comes with absolutely no warranties whatsoever. By using Scour, you take full responsibility for any and all outcomes that result.
