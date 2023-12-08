import pandas as pd
import numpy as np
import matplotlib.pyplot as plt


segsize_means_1k = []
segsize_means_volvo = []
segsize_means_sepsis = []


file_segsize_1k = {
    '0.1': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_0.1/test_01_1.csv',
    '0.2': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_0.2/test_02_1.csv',
    '0.5': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_0.5/test_05_1.csv',
    '1':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_1/test_1_1.csv',
    '2':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_2/test_2_1.csv',
    '3':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_3/test_3_1.csv',
    '4':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_4/test_4_1.csv',
    '4.5': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_4.5/test_45_1.csv',
    '5':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_5/test_5_1.csv',
    '6':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_6/test_6_1.csv',
}

file_segsize_volvo = {
    '0.1': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_0.1/test_01_1.csv',
    '0.2': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_0.2/test_02_1.csv',
    '0.5': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_0.5/test_05_1.csv',
    '1': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_1/test_1_1.csv',
    '2': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_2/test_2_1.csv',
    '3': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_3/test_3_1.csv',
    '4': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_4/test_4_1.csv',
    '4.5': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_4.5/test_45_1.csv',
    '5': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_5/test_5_1.csv',
    '6': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_6/test_6_1.csv',

}

file_segsize_sepsis = {
    '0.1': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_0.1/test_01_1.csv',
    '0.2': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_0.2/test_02_1.csv',
    '0.5': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_0.5/test_05_1.csv',
    '1': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_1/test_1_1.csv',
    '2': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_2/test_2_1.csv',
    '3': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_3/test_3_1.csv',
    '4': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_4/test_4_1.csv',
    '4.5': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_4.5/test_45_1.csv',
    '5': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_5/test_5_1.csv',
    '6': '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_6/test_6_1.csv',

}

# DATASET LOG 1K

#Read CSV and means second column
for nome, file in file_segsize_1k.items():
    df_test = pd.read_csv(file, decimal='.', header=1)
    mean_column = df_test.iloc[:, 1].mean()

    #Convert Means from KB to MB
    mean_mb = mean_column * 0.000001
    segsize_means_1k.append(mean_mb)


segsize_means_1k.insert(0,0)
list_x = list(file_segsize_volvo.keys())
list_x.insert(0,0)
#name_xaxis_1k = [float(x) for x in list_x]

name_xaxis_1k = [float(x) for x in list_x]




# DATASET LOG VOLVO

for nome, file in file_segsize_volvo.items():
    df_test_volvo = pd.read_csv(file, decimal='.', header=1)
    mean_column_volvo = df_test_volvo.iloc[:, 1].mean()

    #Convert Means from KB to MB
    mean_mb_volvo = mean_column_volvo * 0.000001
    segsize_means_volvo.append(mean_mb_volvo)

segsize_means_volvo.insert(0,0)
list_x = list(file_segsize_volvo.keys())
list_x.insert(0,0)
name_xaxis_volvo = [float(x) for x in list_x]




# DATASET LOG SEPSIS

for nome, file in file_segsize_sepsis.items():
    df_test_sepsis = pd.read_csv(file, decimal='.', header=1)
    mean_column_sepsis = df_test_sepsis.iloc[:, 1].mean()

    #Convert Means from KB to MB
    mean_mb_sepsis = mean_column_sepsis * 0.000001
    segsize_means_sepsis.append(mean_mb_sepsis)


segsize_means_sepsis.insert(0,0)
list_x = list(file_segsize_volvo.keys())
list_x.insert(0,0)
name_xaxis_sepsis = [float(x) for x in list_x]



print(max(segsize_means_1k))


# PLOT

plt.style.use("seaborn-v0_8-bright")
plt.figure(figsize=(16, 9))

# Create a line plot for the dataset motivating scenario
plt.plot(name_xaxis_1k, segsize_means_1k, label='Motivating scenario', color='skyblue', linewidth=5, marker = 'p', markersize=10, alpha= 1)

# Create a line plot for the dataset volvo
plt.plot(name_xaxis_volvo, segsize_means_volvo, label='BPIC 2013', color='tomato', linewidth=5, marker = 'v', markersize=10, alpha= 0.9)

# Create a line plot for the dataset sepsis
plt.plot(name_xaxis_sepsis, segsize_means_sepsis,  label='Sepsis', color='mediumaquamarine', linewidth=5, marker = '*', markersize=10, alpha= 1)


# STABILIZATION LINES

plt.vlines(x = 4.5, ymin = 0, ymax = 23.8, color = "skyblue", linestyle='dashed', linewidth = 4, alpha = 0.8, label = 'Stabilization point')
plt.vlines(x = 3, ymin = 0, ymax = 15.5, color = "tomato", linestyle='dashed', linewidth = 4, alpha = 0.8, label = 'Stabilization point')
plt.vlines(x = 4, ymin = 0, ymax = 19, color = "mediumaquamarine", linestyle='dashed', linewidth = 4, alpha = 0.8, label = 'Stabilization point')


plt.xlabel('Segment size (MB)', fontsize = 22, labelpad= 15)
plt.ylabel('Memory usage (MB)', fontsize = 22,  labelpad= 15)
plt.xticks(fontsize=22)
plt.yticks(fontsize=22)
plt.legend(['Motivating scenario', 'BPIC_2013', 'Sepsis', 'Stabilization point'],  loc='upper left', fontsize=18)
plt.grid(True, linestyle='--')

ax = plt.gca()
leg = ax.get_legend()
leg.legendHandles[3].set_color('grey')

plt.xlim([0, 6.4])
plt.ylim([0, 25])

plt.tight_layout()
plt.savefig('lineplot_segsize_syth_combined.pdf')


plt.savefig('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_memoryusage/memoryusage4.pdf')
plt.show()
exit()