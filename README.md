# CONFINE: Preserving data secrecy in inter-organizational process mining
The CONFINE framework represents a decentralized paradigm for inter-organizational process mining, utilizing trusted applications operating within Trusted Execution Environments (TEEs) to uphold the tenets of data privacy and confidentiality. This repository houses the prototype implementation of the principal components constituting the framework, namely:

- Provisioners: HTTP servers responsible for delivering event logs designated for mining.
- Secure Miners: [EGo](https://www.edgeless.systems/products/ego/) Intel SGX trusted applications retrieving and merging event logs to be fed into process mining algorithms.

## Framework overview
## Repository
The repository is structured as follows:
-  [/src/](https://github.com/Process-in-Chains/CONFINE/tree/main/src): the root folder of the implementation source code
    - [/src/secure-miner/](https://github.com/Process-in-Chains/CONFINE/tree/main/src/secure-miner) contains the Secure Miner implementation as an EGo Intel SGX application
    - [/src/provisioner/](https://github.com/Process-in-Chains/CONFINE/tree/main/src/provisioner) contains the Log Provider component implementation in GO
-  [/evaluation/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation): folder containing datasets and results of our tests
    - [/evaluation/convergence/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/convergence) contains the convergence test data 
    - [/evaluation/memoryusage/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage) includes the memory usage tests data 
    - [/evaluation/scalability/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability) contains the data of the scalability tests
## Installation
## Setup and run
### Provisioner node
### Secure miner
## Evaluation
### Datasets
### Tests
