import pandas as pd
#import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime



# Read CSV

df = pd.read_csv(
    '/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_memoryusage/test_nocomput/test_1_new.csv', decimal='.', header=0)

# SECONDS
df['Timestamp'] = df['Timestamp'].apply(lambda x: datetime.utcfromtimestamp(x))


# Calculate first boot timestamp
start_time = df['Timestamp'].min()


"""# Transform timestamps into minutes and seconds
df_2['Durata (Minuti)'] = (df_2['Timestamp'] - start_time).dt.total_seconds() / 60
df_2['Durata (Minuti)'] = df_2['Durata (Minuti)'].apply(lambda x: '{:02}:{:02}'.format(int(x), int((x % 1) * 60)))
"""

# Transform timestamps into seconds
df['Durata (Secondi)'] = (df['Timestamp'] - start_time).dt.total_seconds()

# Calculate total runtime
total_runtime_seconds = df['Durata (Secondi)'].max() - df['Durata (Secondi)'].min()


# Normalize 'Durata (Secondi)'
df['Durata Normalizzata'] = (df['Durata (Secondi)'] - df['Durata (Secondi)'].min()) / total_runtime_seconds
df['Durata Normalizzata'] = df['Durata Normalizzata'] * 48

# Print the normalized values on a 0-100 scale
# print(df['Durata Normalizzata'])


# Convert Bytes in MegaBytes
df['Memory usage (MB)'] = df['RAM Usage (Bytes)'] / 1048576



# Unify the dataset
result = df.groupby('Durata Normalizzata')['Memory usage (MB)'].mean().reset_index()
pd.options.display.float_format = '{:.2f}'.format
# print(result)


"""
# Add the latest timestamp to the graph
ultimo_timestamp = df_2['Durata Normalizzata'].max()
selected_timestamps = result['Durata Normalizzata'][::5].tolist()
selected_timestamps.append(ultimo_timestamp)
"""

# LINES

"""
TESTMODE - INITIALIZATION STARTED AT: 1700497715461
TESTMODE - FIRST ATTESTATION AT: 1700497724188
TESTMODE - FIRST SEGMENT RECEIVED AT: 1700497725625
"""

inizialization = 1700497715461
first_segment_recieved =  1700497725625
attestation = 1700497724188



diff_seconds_segment_received = (first_segment_recieved - inizialization)/1000
diff_seconds_norm = (diff_seconds_segment_received - df['Durata (Secondi)'].min()) / total_runtime_seconds
diff_seconds_norm = diff_seconds_norm * 48
print(diff_seconds_norm)


diff_seconds_att = (attestation - inizialization)/1000
diff_seconds_norm_att = ( diff_seconds_att - df['Durata (Secondi)'].min()) / total_runtime_seconds
diff_seconds_norm_att = diff_seconds_norm_att * 48
print(diff_seconds_norm_att)


# PLOT
plt.style.use("seaborn-v0_8-bright")
plt.figure(figsize=(16,9))

graph = plt.plot(result['Durata Normalizzata'],result['Memory usage (MB)'], color = 'purple',linewidth=3)
plt.xticks(fontsize=22)
plt.yticks(fontsize=22)
plt.xlabel('Run completion percentage', fontsize = 22, labelpad= 15)
plt.ylabel('Memory usage (MB)', fontsize = 22,  labelpad= 15)
plt.grid(True, linestyle='--')
plt.tight_layout()

plt.xlim([0, 48])
plt.ylim([0, df['Memory usage (MB)'].max()])


plt.axvline(x = int(diff_seconds_norm_att), color = "darkslategray", linestyle = "dashed", linewidth = 4, label = 'First attestation', alpha = 0.8)
plt.axvline(x = int(diff_seconds_norm), color = "firebrick", linestyle='dashed', dashes=(2, 2), linewidth = 4, label = 'First segment recieved', alpha = 0.8)

plt.legend(['Memory usage trend', 'First attestation', 'First segment recieved'], loc='upper right', fontsize=18)
# plt.legend(['Memory usage trend', 'First segment recieved', 'First computation', 'First attestation'], loc='upper right', fontsize=18, framealpha=1, edgecolor="black", fancybox=False)

plt.fill_between(result['Durata Normalizzata'],result['Memory usage (MB)'], color = 'orchid', alpha = 0.1)
plt.tight_layout()
plt.savefig('/Users/luca/Documents/PythonProjects/TEE_Evaluation/test_memoryusage/memoryusage1.pdf')
plt.show()
exit()


