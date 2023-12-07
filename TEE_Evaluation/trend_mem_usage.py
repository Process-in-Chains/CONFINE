import pandas as pd
#import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime


segsize_means = []

"""
file_segsize = {
    '4':'/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_4/test_4.csv',
}"""

# Read CSV

df_2 = pd.read_csv(
    '/Users/luca/Documents/PythonProjects/TEE_Evaluation/segment_size/test_simulation_1k/segsize_4/test_4.csv', decimal='.', header=0)
# df_2['Timestamp'] = df_2['Timestamp'].apply(lambda x: datetime.utcfromtimestamp(x).strftime('%Y-%m-%d %H:%M:%S'))
df_2['Timestamp'] = df_2['Timestamp'].apply(lambda x: datetime.utcfromtimestamp(x))

# Calculate first boot timestamp
start_time = df_2['Timestamp'].min()

# Transform timestamps into minutes and seconds
df_2['Durata (Minuti)'] = (df_2['Timestamp'] - start_time).dt.total_seconds() / 60
df_2['Durata (Minuti)'] = df_2['Durata (Minuti)'].apply(lambda x: '{:02}:{:02}'.format(int(x), int((x % 1) * 60)))

# Convert Bytes in MegaBytes
df_2['RAM Usage (MB)'] = df_2['RAM Usage (Bytes)'] / 1048576


# Unify the dataset
result = df_2.groupby('Durata (Minuti)')['RAM Usage (MB)'].mean().reset_index()
pd.options.display.float_format = '{:.2f}'.format
print(result)

# Add the latest timestamp to the graph

ultimo_timestamp = df_2['Durata (Minuti)'].max()
selected_timestamps = result['Durata (Minuti)'][::5].tolist()
selected_timestamps.append(ultimo_timestamp)



# PLOT

plt.style.use("seaborn-v0_8-bright")
plt.figure(figsize=(16,9))
graph = plt.plot(result['Durata (Minuti)'],result['RAM Usage (MB)'], color = 'darkslateblue',linewidth=3, marker='.')
plt.xticks(selected_timestamps, fontsize=15)
plt.yticks(fontsize=12)
plt.xlabel('Timestamp')
plt.ylabel('RAM Usage (MB)')
plt.grid(True, linestyle='--')
plt.tight_layout()
# leg = plt.legend (loc='lower right', framealpha=1, edgecolor="black", fancybox=False)

plt.title("RAM USAGE for single run")
plt.savefig('Ram_usage_per_TS.pdf')
plt.show()
exit()
