import pandas as pd
#import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime



# Read CSV
#df_simulation = pd.read_csv('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_simulation_1k/segsize_2/test_2_1.csv', decimal='.', header=0,)
df_simulation = pd.read_csv('./declareConformanceConfiguration/incrementalDeclare.segsize100.motivating.csv', decimal='.', header=0)

#df_sepsis = pd.read_csv('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_sepsis/segsize_2/test_2_1.csv', decimal='.', header=0,)
df_sepsis = pd.read_csv('./declareConformanceConfiguration/incrementalDeclare.segsize100.sepsis.csv', decimal='.', header=0,)

#df_volvo = pd.read_csv('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_segment_size/test_volvo/segsize_2/test_2_1.csv', decimal='.', header=0,)
df_volvo = pd.read_csv('./declareConformanceConfiguration/incrementalDeclare.segsize100.volvo.csv', decimal='.', header=0,)


# Convert in datetime
df_simulation['Timestamp'] = df_simulation['Timestamp'].apply(lambda x: datetime.utcfromtimestamp(x/1000))
df_sepsis['Timestamp'] = df_sepsis['Timestamp'].apply(lambda x: datetime.utcfromtimestamp(x/1000))
df_volvo['Timestamp'] = df_volvo['Timestamp'].apply(lambda x: datetime.utcfromtimestamp(x/1000))

df_volvo = df_volvo.iloc[::2]
df_sepsis = df_sepsis.iloc[::1]

# Calculate first boot timestamp
start_time = df_simulation['Timestamp'].min()
start_time_sepsis = df_sepsis['Timestamp'].min()
start_time_volvo = df_volvo['Timestamp'].min()

# Transform timestamps into seconds
df_simulation['Durata (Seconds)'] = (df_simulation['Timestamp'] - start_time).dt.total_seconds()
df_sepsis['Durata (Seconds)'] = (df_sepsis['Timestamp'] - start_time_sepsis).dt.total_seconds()
df_volvo['Durata (Seconds)'] = (df_volvo['Timestamp'] - start_time_volvo).dt.total_seconds()


# Calculate total runtime
total_runtime_seconds_simulation = df_simulation['Durata (Seconds)'].max() - df_simulation['Durata (Seconds)'].min()
total_runtime_seconds_sepsis = df_sepsis['Durata (Seconds)'].max() - df_sepsis['Durata (Seconds)'].min()
total_runtime_seconds_volvo = df_volvo['Durata (Seconds)'].max() - df_volvo['Durata (Seconds)'].min()

# Normalize 'Durata (Secondi)' Simulation
df_simulation['Durata Normalizzata'] = (df_simulation['Durata (Seconds)'] - df_simulation['Durata (Seconds)'].min()) / total_runtime_seconds_simulation
df_simulation['Durata Normalizzata'] = df_simulation['Durata Normalizzata'] * 100
# Normalize 'Durata (Secondi)' Sepsis
df_sepsis['Durata Normalizzata'] = (df_sepsis['Durata (Seconds)'] - df_sepsis['Durata (Seconds)'].min()) / total_runtime_seconds_sepsis
df_sepsis['Durata Normalizzata'] = df_sepsis['Durata Normalizzata'] * 100
# Normalize 'Durata (Secondi)' Volvo
df_volvo['Durata Normalizzata'] = (df_volvo['Durata (Seconds)'] - df_volvo['Durata (Seconds)'].min()) / total_runtime_seconds_volvo
df_volvo['Durata Normalizzata'] = df_volvo['Durata Normalizzata'] * 100



# Convert Bytes in MegaBytes
df_simulation['Memory usage (MB)'] = df_simulation['RAM Usage (Bytes)'] / 1048576
df_sepsis['Memory usage (MB)'] = df_sepsis['RAM Usage (Bytes)'] / 1048576
df_volvo['Memory usage (MB)'] = df_volvo['RAM Usage (Bytes)'] / 1048576

# Unify the dataset
result = df_simulation.groupby('Durata Normalizzata')['Memory usage (MB)'].mean().reset_index()
result_sepsis = df_sepsis.groupby('Durata Normalizzata')['Memory usage (MB)'].mean().reset_index()
result_volvo = df_volvo.groupby('Durata Normalizzata')['Memory usage (MB)'].mean().reset_index()

pd.options.display.float_format = '{:.2f}'.format
# print(result)

"""# Add the latest timestamp to the graph
ultimo_timestamp = df_simulation['Durata (Minuti)'].max()
selected_timestamps = result['Durata (Minuti)'][::5].tolist()
selected_timestamps.append(ultimo_timestamp)

ultimo_timestamp_sepsis = df_sepsis['Durata (Minuti)'].max()
selected_timestamps_sepsis = result_sepsis['Durata (Minuti)'][::5].tolist()
selected_timestamps_sepsis.append(ultimo_timestamp_sepsis)


ultimo_timestamp_volvo = df_volvo['Durata (Minuti)'].max()
selected_timestamps_volvo = result_volvo['Durata (Minuti)'][::5].tolist()
selected_timestamps_volvo.append(ultimo_timestamp_volvo)"""



"""selected_timestamps = result_volvo['Durata (Minuti)'][result_volvo.index % 5 == 0]

# Estrai i timestamp di ogni 5 secondi per i dataframe sepsis e volvo
selected_timestamps_sepsis = result_sepsis['Durata (Minuti)'][result_sepsis.index % 5 == 0]
# selected_timestamps_volvo = result_volvo['Durata (Minuti)'][result_volvo.index % 5 == 0]
selected_timestamps_simulation = result['Durata (Minuti)'][result.index % 5 == 0]

"""



# PLOT
plt.style.use("seaborn-v0_8-bright")
plt.figure(figsize=(16,9))

plt.plot(result['Durata Normalizzata'], result['Memory usage (MB)'], label='Motivating scenario', color='deepskyblue', linewidth=4, marker = 'p', markersize=8)

# Create a line plot for the dataset volvo
plt.plot(result_volvo['Durata Normalizzata'], result_volvo['Memory usage (MB)'], label='BPIC 2013', color='tomato', linewidth=4, marker = 'v', markersize=8)

# Create a line plot for the dataset sepsis
plt.plot(result_sepsis['Durata Normalizzata'], result_sepsis['Memory usage (MB)'], label='Sepsis', color='forestgreen', linewidth=4, marker='*', markersize=8)

plt.xticks(fontsize=30)
plt.yticks(fontsize=30)
plt.xlabel('Run completion percentage', fontsize = 30, labelpad= 15)
plt.ylabel('Memory usage (MB)', fontsize = 30,  labelpad= 15)
plt.grid(True, linestyle='--')
plt.tight_layout()

plt.xlim([0, 101])
plt.ylim([0,20])


plt.legend (loc='upper right', fontsize=25)

#plt.fill_between(result['Durata Normalizzata'],result['Memory usage (MB)'], color = 'azure')
plt.tight_layout()
#plt.savefig('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_memoryusage/memoryusage3.pdf')
plt.savefig('multilog_declareConformanceConfiguration.pdf')
plt.show()
exit()
