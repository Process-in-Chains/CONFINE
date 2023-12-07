import pandas as pd
#import numpy as np
import matplotlib.pyplot as plt


segsize_means = []

file_segsize = {
    '0.1':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_0.1/test_01.csv',
    '0.2':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_0.2/test_02.csv',
    '0.5':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_0.5/test_05.csv',
    '1':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_1/test_1.csv',
    '1.5':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_1.5/test_15.csv',
    '2':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_2/test_2.csv',
    '3':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_3/test_3.csv',
    '4':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_4/test_4.1.csv',
    '4.5':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_4.5/test_45.csv',
    '5':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_5/test_5.csv',
    '6':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_6/test_6.csv',
}

#SINGLE CSV
"""df_2 = pd.read_csv(
    '/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation/segsize_2/test_2.csv', decimal='.', header=1)
#print(df_2)
mean_segsize_2 = df_2.iloc[:, 1].mean()
print(mean_segsize_2)
exit()"""


#Read CSV and means second column
for nome, file in file_segsize.items():
    df_test = pd.read_csv(file, decimal='.', header=1)
    mean_column = df_test.iloc[:, 1].mean()

    #Convert Means from KB to MB
    mean_mb = mean_column * 0.000001
    segsize_means.append(mean_mb)


name_xaxis = list(file_segsize.keys())



# PLOT

plt.style.use("seaborn-v0_8-bright")
plt.figure(figsize=(16,9))
graph = plt.bar(name_xaxis, segsize_means, color = 'lightblue')
plt.xlabel('Segment Size (MB)')
plt.ylabel('RAM Usage (MB)')
# plt.xticks(fontsize=15)
# plt.yticks(fontsize=12)
plt.tight_layout()

plt.title("RAM USAGE for different segment size")
plt.savefig('barplot_segsize.pdf')
plt.show()
exit()


