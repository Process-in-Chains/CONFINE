import pandas as pd
#import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime



# Read CSV

#df = pd.read_csv('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_memoryusage/test_new/test_1_1.csv', decimal='.', header=0)
#df = pd.read_csv('/content/drive/MyDrive/TEEs_for_Inter_Organization_Process_Mining/tests/test_memory_usage/test_full_run/test_1_1.csv', decimal='.', header=0)


#df = pd.read_csv('/content/drive/MyDrive/TEEs_for_Inter_Organization_Process_Mining/tests/Journal_tests/runtime_test/motivating_segsize2000.csv', decimal='.', header=0)

#conformance
#df = pd.read_csv('/content/drive/MyDrive/TEEs_for_Inter_Organization_Process_Mining/tests/Journal_tests/runtime_test/conformance_checking/declareConformance_segsize2000.csv', decimal='.', header=0)
df = pd.read_csv('./incremental/incrementalHeuristics.segsize100.motivating.csv', decimal='.', header=0)
#df = pd.read_csv('/content/drive/MyDrive/TEEs_for_Inter_Organization_Process_Mining/tests/Journal_tests/runtime_test/classicHeuristics.segsize100.motivating.csv', decimal='.', header=0)
# SECONDS
df['Timestamp'] = df['Timestamp'].apply(lambda x: datetime.utcfromtimestamp(x/1000))

# Calculate first boot timestamp
start_time = df['Timestamp'].min()
# Transform timestamps into seconds
df['Durata (Secondi)'] = (df['Timestamp'] - start_time).dt.total_seconds()
# Calculate total runtime
total_runtime_seconds = df['Durata (Secondi)'].max() - df['Durata (Secondi)'].min()


# Normalize 'Durata (Secondi)'
df['Durata Normalizzata'] = (df['Durata (Secondi)'] - df['Durata (Secondi)'].min()) / total_runtime_seconds
df['Durata Normalizzata'] = df['Durata Normalizzata'] * 100



# Convert Bytes in MegaBytes
df['Memory usage (MB)'] = df['RAM Usage (Bytes)'] / 1048576

# Unify the dataset
result = df.groupby('Durata Normalizzata')['Memory usage (MB)'].mean().reset_index()
pd.options.display.float_format = '{:.2f}'.format
# print(result)



# LINES

"""TESTMODE - TEST STARTED AT:  1700495054658
TESTMODE - INITIALIZATION STARTED AT: 1700495054660
TESTMODE - FIRST ATTESTATION AT: 1700495063742
TESTMODE - FIRST SEGMENT RECEIVED AT: 1700495065158
TESTMODE - FIRST COMPUTATION AT: 1700495068178
TESTMODE - TEST ENDED AT:  1700495091025"""

inizialization = 1719570436055
first_segment_recieved =  1719570436683
first_computation = 1719570436983
attestation = 1719570436583

#inizialization = 1719570593068
#first_segment_recieved =  1719570593708
#first_computation = 1719570600248
#attestation = 1719570593600



diff_seconds_segment_received = (first_segment_recieved - inizialization)/1000
diff_seconds_norm = (diff_seconds_segment_received - df['Durata (Secondi)'].min()) / total_runtime_seconds
diff_seconds_norm = diff_seconds_norm * 100
print(diff_seconds_norm)

diff_seconds_comp = (first_computation - inizialization)/1000
diff_seconds_norm_comp = ( diff_seconds_comp - df['Durata (Secondi)'].min()) / total_runtime_seconds
diff_seconds_norm_comp = diff_seconds_norm_comp * 100
print(diff_seconds_norm_comp)


diff_seconds_att = (attestation - inizialization)/1000
diff_seconds_norm_att = ( diff_seconds_att - df['Durata (Secondi)'].min()) / total_runtime_seconds
diff_seconds_norm_att = diff_seconds_norm_att * 100
print(diff_seconds_norm_att)


# PLOT
plt.style.use("seaborn-v0_8-bright")
plt.figure(figsize=(16,9))

#plt.plot(result['Durata Normalizzata'],result['Memory usage (MB)'], color = 'purple',linewidth=4, marker='.')
plt.plot(result['Durata Normalizzata'],result['Memory usage (MB)'], color = 'steelblue',linewidth=4, marker='.')
# plt.plot(result_2['Durata Normalizzata'],result_2['Memory usage (MB)'].fillna(0), color = 'purple',linewidth=3, marker='.')


plt.xticks(fontsize=30)
plt.yticks(fontsize=30)
plt.xlabel('Run completion percentage', fontsize = 30, labelpad= 15)
plt.ylabel('Memory usage (MB)', fontsize = 30,  labelpad= 15)
plt.grid(True, linestyle='--')
plt.tight_layout()

plt.xlim([0, 101])
#plt.ylim([0, df['Memory usage (MB)'].max()])
plt.ylim([0, 25])


plt.axvline(x = int(diff_seconds_norm_att), color = "teal", linestyle='dashed', linewidth = 5, label = 'First attestation', alpha = 0.8)
plt.axvline(x = int(diff_seconds_norm), color = "crimson", linestyle='dashed', dashes=(2, 2), linewidth = 5, label = 'First segment recieved', alpha = 0.8)
plt.axvline(x = int(diff_seconds_norm_comp), color = "goldenrod", linestyle='dotted', dashes = (1,1), linewidth = 5, label = 'First computation', alpha = 0.8)


plt.legend(['Memory usage trend', 'First attestation', 'First segment recieved', 'First computation'], loc='upper left', fontsize=25, framealpha=1)
# plt.legend(['Memory usage trend', 'First segment recieved', 'First computation', 'First attestation'], loc='upper right', fontsize=18, framealpha=1, edgecolor="black", fancybox=False)

plt.fill_between(result['Durata Normalizzata'],result['Memory usage (MB)'], color = 'azure')
#plt.fill_between(result['Durata Normalizzata'],result['Memory usage (MB)'], color = 'orchid', alpha = 0.1)
plt.tight_layout()
#plt.savefig('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_memoryusage/memoryusage2.pdf')
plt.savefig('plot_heuriticsminerIncremental.pdf')
#plt.show()
exit()