# CONFINE: Preserving data secrecy in inter-organizational process mining
The CONFINE framework represents a decentralized paradigm for inter-organizational process mining, utilizing trusted applications operating within Trusted Execution Environments (TEEs) to uphold the tenets of data privacy and confidentiality. This repository houses the prototype implementation of the principal components constituting the framework, namely:

- Provisioners: HTTP servers responsible for delivering event logs designated for mining.
- Secure Miners: [EGo](https://www.edgeless.systems/products/ego/) Intel SGX trusted applications retrieving and merging event logs to be fed into process mining algorithms.

## Framework overview
![deploymentdiagram-1](https://github.com/Process-in-Chains/CONFINE/assets/60829979/5c8dded3-5f04-42a7-a9d0-1a4583ddf708)

Our framework involves different information systems running on multiple machines. An organization can take at least one of the following roles: provisioning if it delivers local event logs to be collaboratively mined; mining if it applies process mining algorithms using event logs retrieved from provisioners. Depending on the played role, nodes come endowed with a Provisioner or a Secure Miner component, or both. 
Provisioner Nodes host the Provisioner's components, encompassing the Log Recorder and the Log Provider. 
The Miner Node is characterized by two distinct execution environments: the Operating System(OS) and the Trusted Execution Environment (TEE). TEEs establish isolated contexts separate from the OS, safeguarding code and data through hardware-based encryption mechanisms. We leverage the security guarantees provided by TEEs to protect a Trusted App responsible for fulfilling the functions of the Secure Miner and its associated sub-components. 
The Secure Miner exchange messages with Provisioners according to the CONFINE protocol. After the proper execution of the CONFINE protocol, the trusted app implementing the Secure Miner is able to retrieve event logs, merge them and elaborate their aggregation in the Trusted Execution Environment.


## Repository
The main content of the repository is structured as follows:
-  [/src/](https://github.com/Process-in-Chains/CONFINE/tree/main/src): the root folder of the implementation source code
    - [/src/secure-miner/](https://github.com/Process-in-Chains/CONFINE/tree/main/src/secure-miner) contains the Secure Miner implementation as an EGo Intel SGX application
    - [/src/provisioner/](https://github.com/Process-in-Chains/CONFINE/tree/main/src/provisioner) contains the Log Provider component implementation in GO
    - [/src/mining-data/](https://github.com/Process-in-Chains/CONFINE/tree/main/src/mining-data) contains the metadata required for the execution of the CONFINE protocol
-  [/evaluation/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation): folder containing datasets and results of our tests
    - [/evaluation/convergence/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/convergence) contains the convergence test data 
    - [/evaluation/memoryusage/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage) includes the memory usage tests data 
    - [/evaluation/scalability/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability) contains the data of the scalability tests
## Installation
## Setup and run
### Provisioner node
### Secure miner
## Evaluation
The following section contains the experimental toolbench used to evaluate the effectiveness of CONFINE, presented in the paper Section 6. Evaluation files can be found in [/evaluation/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation). We conduct convergence analysis to demonstrate the correctness of the collaborative data exchange process. Moreover, we gauge the memory usage with synthetic and real-life event logs, to observe the trend during the enactment of our protocol and assess scalability. 

##### Requirements
To run our Python scripts, the following libraries are required: `os`, `pandas`, `numpy`, `matplotlib`, `scipy`, `sklearn`, `datetime`.

### Tests
We begin with a convergence analysis to demonstrate the correctness of the collaborative data exchange process. We gauge the memory usage with synthetic and real-life event logs, to observe the trend during the enactment of our protocol and assess scalability. 

##### Convergence

To experimentally validate the correctness of our approach in the transmission and computation phases, we run a [/convergence/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/convergence) test. To this end, we created a synthetic event log (available in [/event_log/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/convergence/event_log)) consisting of 1000 cases of 14 events on average by simulating the inter-organizational process of our motivating scenario and we partitioned it in three sub-logs (Respectively Hospital, Specialized clinic and Pharmaceutical company event logs). We run the stand-alone HeuristicsMiner on the former, and processed the latter through our CONFINE toolchain. The convergence results are available in [/output/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/convergence/output) in the form of a workflow net.


##### Memory Usage
[/memory usage/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage)
- (I) No Comput 
- (II) Comput
- (III) Multiple log lines
- (IV) Multiple log segsize

[/event_log/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage/event_log)
[/bpic13/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage/event_log/bpic13)
[/motivating_scenatio/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage/event_log/motivating_scenario)
[/sepsis/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage/event_log/sepsis)




##### Scalability
We examine the [/scalability/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability) of the Secure Miner, focusing on its capacity to efficiently manage an increasing workload in the presence of limited memory resources. We implemented three distinct test configurations gauging runtime memory usage as variations of our motivating scenario log.

- (I) The maximum number of events per case 
- (II) The number of cases
- (III) The number of provisioning organizations
