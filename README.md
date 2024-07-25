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

## Screencast
As follows you can find a screencast that shows how to set up and run the necessary components

https://github.com/user-attachments/assets/66b971bf-a0a8-4e60-97ee-e00796418a43

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
## Installation with Docker
To run the code, we prepared a CONFINE Docker image containing all the necessary requirements in order to run the framework. In order to use it, you need to have [Docker] installed, you can do that by following [this page](https://docs.docker.com/get-docker/). After installation, you can start creating the container to run the framework. After installation, you can pull the images executing the following command.
```
docker pull valeriogoretti9/confine:latest
```
Once the image has been downloaded, you can execute the following command to create a container with the CONFINE image. Before execution, you must be careful to enter your volume path instead of the <INSERT HERE YOUR VOLUME PATH> tag. You can also set a container name by replacing the tag <INSERT HERE YOUR DOCKER CONTAINER NAME>.
```
docker run --volume /var/run/aesmd:/var/run/aesmd -v <INSERT YOUR VOLUME PATH>:/volume --name <INSERT YOUR DOCKER CONTAINER NAME> -ti valeriogoretti9/confine:latest
```
Once the Docker container is created, the following commands allow you to start it.
```
docker start <INSERT YOUR DOCKER CONTAINER NAME>
docker attach <INSERT YOUR DOCKER CONTAINER NAME>
```

## Setup and run
Once the docker container is running, you can proceed with the CONFINE setup. Execute the following commands to clone the CONFINE code and access it.
```
cd volume
apt-get update
apt-get install git
git clone https://github.com/Process-in-Chains/CONFINE.git
cd CONFINE/
```
In the next parts, you can see how to run [provisioner](#provisioner) and [secure miner](#secure-miner) components. 
### Provisioner
At this point, the container is running. 

By executing the following command, you enter the folder dedicated to the provisioner
```
cd mining-data/provision-data/process-01/
```
You have to put the log (in xes format) you want to provide to the Secure Miners into this folder. We already provide different inter-organizational log samples that you can find in this folder.

After that, navigate to /src/provision-data and modify the minerList.json file
```
cd ..
nano minerList.json
```
Append to this file the TLS certificate string of the Secure Miners you want to be accepted by the provisioner. You will see in the Secure Miner section how to get this information.

Now you are ready to run the provisioner. To facilitate its start-up, we prepared the shell script runLogServer.sh in the src folder. Let's navigate there
```
cd ../..
```
and run the runLogServer.sh shell script  
```
./runLogServer.sh -port 8089 -log testing_logs/motivating/pharma.xes -mergekey concept:name -measurement `ego uniqueid app` -skipattestation true
```
with parameters :
- port: the port on which the log server listens for new requests from the Secure Miner. The default value is 8089.
- log: the path of the XES event log in the /src/mining-data/provision-data folder. The default value is 'testing_logs/motivating/pharma.xes'.
- mergekey: the name of the case identifier attribute inside the provided event log. The default value is 'concept:name'
- measurement: the value that identifies the Secure Miner's source code for the remote attestation. The default value 'ego uniqueid app' uses an EGo command to compute this information using the Secure Miner's source code.
- skipattestation: if is set to true the remote attestation phase of the CONFINE protocol will be skipped. The default value is true. **If the Secure Miner is running in simulation mode, this must be set to true**.
### Secure miner
In order to enable communication with other parties involved in the Inter-organizational environment, you have to set your collaborator. To do so, you must edit the file refeece.json. The following commands describe how to do this. To do so, you must edit the file reference.json. First, from the src folder, navigate to the file:
```
cd mining-data/collaborators/process-01/
```
Inside this folder is the references.json file. Edit this file adding your collaborators. The file contains a list of contributors and for each one you need to set the public key and the reference port on which you want to communicate. The following json is an example for an environment with 2 other collaborators.
```
[
  {
    "public_key": "...""
    "http_reference": "..."
  },
  {
    "public_key": "...""
    "http_reference": "..."
  }
]
```
You have to put, inside this direcroty, the log (in xes format) you want to provide in the inter-organizational context into this folder. The name of this xes file must be "event_log.xes". 
Afterwards, you can build the secure miner code and execute it with the following commands.
```
ego-go build -buildvcs=false
ego sign
```
And finally run the secure miner with the command:
```
OE_SIMULATION=1 ego run ./app -segsize 2000 -port 8080 -test true
```
## Evaluation
The following section contains the experimental toolbench used to evaluate the effectiveness of CONFINE, presented in the paper Section 6. Evaluation files can be found in [/evaluation/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation). We conduct convergence analysis to demonstrate the correctness of the collaborative data exchange process. Moreover, we gauge the memory usage with synthetic and real-life event logs, to observe the trend during the enactment of our protocol and assess scalability. 

##### Requirements
To run our Python scripts, the following libraries are required: `os`, `pandas`, `numpy`, `matplotlib`, `scipy`, `sklearn`, `datetime`.

### Tests
For each test we used the event log of the motivating scenario and additionally, the [BPIC](https://www.tf-pm.org/competitions-awards/bpi-challenge) logs, specifically [Sepsis](https://data.4tu.nl/datasets/33632f3c-5c48-40cf-8d8f-2db57f5a6ce7) and [BPIC 2013](https://data.4tu.nl/datasets/1987a2a6-9f5b-4b14-8d26-ab7056b17929) event logs. We further processed these logs to simulate an inter-organizational scenario. We made specific modifications on the scalability tests event logs to allow the observation of different configurations of number of events per case, number of cases and number of provisioning organizations. To run the test, the specific commands to be execute in the terminal are the following:

###### Miner
`ego-go build -buildvcs=false && ego sign enclave.json && OE_SIMULATION=1 ego run ./app -segsize ${SEGSIZE_VALUE} -port ${PORT_NUMBER} -test true`

###### Provisioner
`go build -o logprovision log_provision/log_provision.go && ./logprovision -port ${PORT_NUMBER} -log ${LOG_PATH}`

#####  Output Convergence

To experimentally validate the correctness of our approach in the transmission and computation phases, we run a [convergence](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/convergence) test. To this end, we created a synthetic event log (available in [/event_log/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/convergence/event_log)) consisting of 1000 cases of 14 events on average by simulating the inter-organizational process of our motivating scenario and we partitioned it in three sub-logs (Respectively Hospital, Specialized clinic and Pharmaceutical company event logs). We run the stand-alone HeuristicsMiner on the former, and processed the latter through our CONFINE toolchain. The convergence results are available in [/output/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/convergence/output) in the form of a workflow net.


##### Memory Usage
To evaluate the runtime memory utilization of our CONFINE implementation, we run a [memory usage](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage) test, split into four different configurations:
- In the first test, we excluded the computation phase by leaving the HeuristicsMiner inactive, so as to isolate execution from mining-specific operations. In this case we set the value of the segment size `${SEGSIZE_VALUE}` to 2000 KiloBytes and we used the same synthetic [event_log](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage/event_log/motivating_scenario) of our motivating scenario. 
- In the second test, on the other hand, we included the computation phase using the same synthetic event log and setting the segment size value as in the first test. The results of the first two tests are available in [/output_test_motivating_scenario/](https://github.com/Process-in-Chains/CONFINE/blob/main/evaluation/memoryusage/output_test_motivating_scenario/test_mem.csv). 
- In the third test, we also gauged the runtime memory usage of two public real-world event logs too. Since those are intra-organizational event logs, we split the contents to mimic an inter-organizational context. In particular, we separated the [Sepsis](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage/event_log/sepsis) event log based on the distinction between normal-care and intensive-care paths, as if they were conducted by two distinct organizations. Similarly, we processed the [BPIC 2013](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage/event_log/bpic13) event log to sort it out into the three departments of the Volvo IT incident management system.
- In the fourth test, we conducted an additional test to examine the trend of memory usage as the segment size varies with all the aforementioned event logs. The results of the test are available in [/output_test_segment_size/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/memoryusage/output_test_segment_size).


##### Scalability
We examine the [scalability](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability) of the Secure Miner, focusing on its capacity to efficiently manage an increasing workload in the presence of limited memory resources. We implemented three distinct test configurations gauging runtime memory usage as variations of our motivating scenario log.

- To conduct the test on the maximum number of events, we modified the motivating scenario [event log](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability/event_log/log_test_maxevents) by adding a loop back from the final to the initial activity of the process model, progressively increasing the number of iterations. The results of the test are available in [/output_test_max_events/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability/output_test_max_events).
- Concerning the test on the number of cases, we simulated additional process instances, building new  [event logs](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability/event_log/log_test_cases).  The results of the test are available in [/output_test_cases/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability/output_test_cases).
- Finally, for the assessment of the number of organizations, the test necessitated the distribution of the process model activities’ into a variable number of pools, each representing a different organization. Event logs are available in [/log_test_organizations/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability/event_log/log_test_organizations), results of the test are available in [/output_test_organizations/](https://github.com/Process-in-Chains/CONFINE/tree/main/evaluation/scalability/output_test_organizations).



[//]: # 
[Docker]: <https://www.docker.com>
