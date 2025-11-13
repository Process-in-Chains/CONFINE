import pandas as pd
import numpy as np
import matplotlib.pyplot as plt


segsize_means_1k = []
segsize_means_volvo = []
segsize_means_sepsis = []


file_segsize_1k = {
    '0.1': './motivating/MS_100.csv',
    '0.2': './motivating/MS_200.csv',
    '0.5': './motivating/MS_500.csv',
    '1':'./motivating/MS_1000.csv',
    '2':'./motivating/MS_2000.csv',
    '3':'./motivating/MS_3000.csv',
    '4':'./motivating/MS_4000.csv',
    '4.5': './motivating/MS_4500.csv',
    '5':'./motivating/MS_5000.csv',
    '6':'./motivating/MS_6000.csv',
    '7':'./motivating/MS_7000.csv',
}

file_segsize_volvo = {
    '0.1': './volvo/V_100.csv',
    '0.2': './volvo/V_200.csv',
    '0.5': './volvo/V_500.csv',
    '1': './volvo/V_1000.csv',
    '2': './volvo/V_2000.csv',
    '3': './volvo/V_3000.csv',
    '4': './volvo/V_4000.csv',
    '4.5': './volvo/V_4500.csv',
    '5': './volvo/V_5000.csv',
    '6': './volvo/V_6000.csv',
    '7': './volvo/V_7000.csv',

}

file_segsize_sepsis = {
    '0.1': './sepsis/S_100.csv',
    '0.2': './sepsis/S_200.csv',
    '0.5': './sepsis/S_500.csv',
    '1': './sepsis/S_1000.csv',
    '2': './sepsis/S_2000.csv',
    '3': './sepsis/S_3000.csv',
    '4': './sepsis/S_4000.csv',
    '4.5': './sepsis/S_4500_2.csv',
    '5': './sepsis/S_5000.csv',
    '6': './sepsis/S_6000.csv',
    '7': './sepsis/S_7000.csv',

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

blue='deepskyblue'
red='tomato'
green='forestgreen'
# Create a line plot for the dataset motivating scenario
plt.plot(name_xaxis_1k, segsize_means_1k, label='Motivating scenario', color=blue, linewidth=5, marker = 'p', markersize=10)

# Create a line plot for the dataset volvo
plt.plot(name_xaxis_volvo, segsize_means_volvo, label='BPIC 2013', color=red, linewidth=5, marker = 'v', markersize=10)

# Create a line plot for the dataset sepsis
plt.plot(name_xaxis_sepsis, segsize_means_sepsis,  label='Sepsis', color=green, linewidth=5, marker = '*', markersize=10)


# STABILIZATION LINES

plt.vlines(x = 4.5, ymin = 0, ymax = 36, color = blue, linestyle='dashed', linewidth = 5, alpha = 1, label = 'Stabilization point')
plt.vlines(x = 3, ymin = 0, ymax = 29, color = red, linestyle='dashed', linewidth = 5, alpha = 1, label = 'Stabilization point')
plt.vlines(x = 4, ymin = 0, ymax = 26, color = green, linestyle='dashed', linewidth = 5, alpha = 1, label = 'Stabilization point')


plt.xlabel('Segment size (MB)', fontsize = 30, labelpad= 15)
plt.ylabel('Memory usage (MB)', fontsize = 30,  labelpad= 15)
plt.xticks(fontsize=30)
plt.yticks(fontsize=30)
plt.legend(['Motivating scenario', 'BPIC 2013', 'Sepsis', 'Stabilization point'],  loc='upper left', fontsize=25)
plt.grid(True, linestyle='--')

ax = plt.gca()
leg = ax.get_legend()
leg.legendHandles[3].set_color('grey')

plt.xlim([0, 7.4])
plt.ylim([0, 40])

plt.tight_layout()
#plt.savefig('lineplot_segsize_syth_combined.pdf')


#plt.savefig('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_memoryusage/memoryusage4.pdf')
plt.savefig('plot_segmentSize_memoryusage.pdf')
plt.show()
exit()