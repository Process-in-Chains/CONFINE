import pandas as pd
import os
import numpy as np
import matplotlib.pyplot as plt
from scipy.optimize import curve_fit
from sklearn.metrics import r2_score
from scipy.stats import spearmanr
from scipy.stats import kendalltau
from sklearn.linear_model import LinearRegression



base_path = './test_number_traces/'

file_sizes = ['128', '256', '512', '1024', '2048', '4096', '8192']
folder_names = ['segsize_100', 'segsize_1000', 'segsize_10000']
num_files_to_average = 3

def leggi_e_calcola_media(base_path, size, folder_name, num_files_to_average):
    path = os.path.join(base_path, size, folder_name)
    dfs = []

    file_names = [f for f in os.listdir(path) if f.endswith('.csv')]
    file_names.sort()


    for file_name in file_names[:num_files_to_average]:
        # Costruisci il percorso completo del file
        file_path = os.path.join(path, file_name)

        # Leggi il file CSV
        df = pd.read_csv(file_path)

        # Aggiungi il DataFrame alla lista
        dfs.append(df)

    # Concatena i DataFrame nella lista in un unico DataFrame
    combined_df = pd.concat(dfs)

    # Calcola la media per ogni colonna
    media = combined_df.mean()

    return media

# Esempio di utilizzo per tutte le combinazioni di dimensione e cartella
risultati = {}

for size in file_sizes:
    for folder_name in folder_names:
        key = f"{size}_{folder_name}"
        risultati[key] = leggi_e_calcola_media(base_path, size, folder_name, num_files_to_average)

# Ora puoi accedere ai risultati come dizionario di DataFrame
# print(risultati['128_segsize_100'])




segsize_means_100 = []
segsize_means_1000 = []
segsize_means_10000= []

file_segsize_100 = {
    '128':(risultati['128_segsize_100']),
    '256':(risultati['256_segsize_100']),
    '512':(risultati['512_segsize_100']),
    '1024':(risultati['1024_segsize_100']),
    '2048':(risultati['2048_segsize_100']),
    '4096':(risultati['4096_segsize_100']),
    '8192':(risultati['8192_segsize_100'])
}

file_segsize_1000 = {
    '128':(risultati['128_segsize_1000']),
    '256':(risultati['256_segsize_1000']),
    '512':(risultati['512_segsize_1000']),
    '1024':(risultati['1024_segsize_1000']),
    '2048':(risultati['2048_segsize_1000']),
    '4096':(risultati['4096_segsize_1000']),
    '8192':(risultati['8192_segsize_1000'])
}


file_segsize_10000 = {
    '128':(risultati['128_segsize_10000']),
    '256':(risultati['256_segsize_10000']),
    '512':(risultati['512_segsize_10000']),
    '1024':(risultati['1024_segsize_10000']),
    '2048':(risultati['2048_segsize_10000']),
    '4096':(risultati['4096_segsize_10000']),
    '8192':(risultati['8192_segsize_10000'])
}



# SEGSIZE 100
# Usa direttamente il DataFrame ottenuto dalla funzione
for size in file_sizes:
    segsize_means_100.append(risultati[f'{size}_segsize_100'].iloc[1] * 0.000001)

segsize_means_100.insert(0, 0)
list_x = list(file_sizes)
list_x.insert(0, 0)
name_xaxis_100 = [float(x) for x in list_x]

df_100 = pd.DataFrame({
    'name_xaxis_100': name_xaxis_100[1:],
    'segsize_means_100': segsize_means_100[1:]
})



# SEGSIZE 1000
# Usa direttamente il DataFrame ottenuto dalla funzione
for size in file_sizes:
    segsize_means_1000.append(risultati[f'{size}_segsize_1000'].iloc[1] * 0.000001)

segsize_means_1000.insert(0, 0)
list_x = list(file_sizes)
list_x.insert(0, 0)
name_xaxis_1000 = [float(x) for x in list_x]


df_1000 = pd.DataFrame({
    'name_xaxis_1000': name_xaxis_1000[1:],
    'segsize_means_1000': segsize_means_1000[1:]
})

# SEGSIZE 10000
# Usa direttamente il DataFrame ottenuto dalla funzione
for size in file_sizes:
    segsize_means_10000.append(risultati[f'{size}_segsize_10000'].iloc[1] * 0.000001)

segsize_means_10000.insert(0, 0)
list_x = list(file_sizes)
list_x.insert(0, 0)
name_xaxis_10000 = [float(x) for x in list_x]

df_10000 = pd.DataFrame({
    'name_xaxis_10000': name_xaxis_10000[1:],
    'segsize_means_10000': segsize_means_10000[1:]
})


"""# CORRELAZIONE LINEARE
print('SEGSIZE 100')
print(df_100['name_xaxis_100'].corr(df_100['segsize_means_100']))
print('SEGSIZE 1000')
print(df_1000['name_xaxis_1000'].corr(df_1000['segsize_means_1000']))
print('SEGSIZE 10000')
print(df_10000['name_xaxis_10000'].corr(df_10000['segsize_means_10000']))
exit()"""

"""#SLOPE

print('SEGSIZE 100')
# Estrai le colonne necessarie
x_values = df_100['name_xaxis_100'].values.astype(float)
y_values = df_100['segsize_means_100'].values.astype(float)
# Calcola la pendenza utilizzando np.polyfit con grado 1 (lineare)
slope_100, intercept_100 = np.polyfit(x_values, y_values, 1)
print(slope_100)
print(intercept_100)
print('SEGSIZE 1000')
# Estrai le colonne necessarie
x_values = df_1000['name_xaxis_1000'].values.astype(float)
y_values = df_1000['segsize_means_1000'].values.astype(float)
# Calcola la pendenza utilizzando np.polyfit con grado 1 (lineare)
slope_1000, intercept_1000 = np.polyfit(x_values, y_values, 1)
print(slope_1000)
print(intercept_1000)
print('SEGSIZE 10000')
# Estrai le colonne necessarie
x_values = df_10000['name_xaxis_10000'].values.astype(float)
y_values = df_10000['segsize_means_10000'].values.astype(float)
# Calcola la pendenza utilizzando np.polyfit con grado 1 (lineare)
slope_10000, intercept_10000 = np.polyfit(x_values, y_values, 1)
print(slope_10000)
print(intercept_10000)
exit()
"""


def relative_ram_scalability(min_ram, max_ram):
    return ((max_ram - min_ram) / min_ram) * 100

def ram_efficiency(relative_ram_scalability, num_nodes):
    return relative_ram_scalability / num_nodes

def ram_speedup(min_ram, max_ram):
    return min_ram / max_ram


"""print('SEGSIZE 100')
# Scalabilità Relativa 100
print(relative_ram_scalability(min(df_100['segsize_means_100']),max(df_100['segsize_means_100'])))
relative_ram_scalability_100 = relative_ram_scalability(min(df_100['segsize_means_100']),max(df_100['segsize_means_100']))
# RAM EFFICIENCY 100
print(ram_efficiency(relative_ram_scalability_100, 7))
# RAM SPEEDUP 100
print(ram_speedup(min(df_100['segsize_means_100']),max(df_100['segsize_means_100'])))

print('SEGSIZE 1000')
# Scalabilità Relativa 1000
print(relative_ram_scalability(min(df_1000['segsize_means_1000']),max(df_1000['segsize_means_1000'])))
relative_ram_scalability_1000 = relative_ram_scalability(min(df_1000['segsize_means_1000']),max(df_1000['segsize_means_1000']))
# RAM EFFICIENCY 1000
print(ram_efficiency(relative_ram_scalability_1000, 7))
# RAM SPEEDUP 1000
print(ram_speedup(min(df_1000['segsize_means_1000']),max(df_1000['segsize_means_1000'])))


print('SEGSIZE 10000')
# Scalabilità Relativa 10000
print(relative_ram_scalability(min(df_10000['segsize_means_10000']),max(df_10000['segsize_means_10000'])))
relative_ram_scalability_10000 = relative_ram_scalability(min(df_10000['segsize_means_10000']),max(df_10000['segsize_means_10000']))
# RAM EFFICIENCY 10000
print(ram_efficiency(relative_ram_scalability_10000, 7))
# RAM SPEEDUP 10000
print(ram_speedup(min(df_10000['segsize_means_10000']),max(df_10000['segsize_means_10000'])))
exit()"""



# PLOT

plt.style.use("seaborn-v0_8-bright")
plt.figure(figsize=(16, 9))


# Create a line plot for the dataset 1K
plt.plot(name_xaxis_100, segsize_means_100, marker='>', markersize = 10,  label='segsize_100', color='navy', linewidth=5)

# Create a line plot for the dataset volvo
plt.plot(name_xaxis_1000, segsize_means_1000, marker='D', markersize = 8, label='segsize_1000', color='mediumorchid', linewidth=5)

# Create a line plot for the dataset sepsis
plt.plot(name_xaxis_10000, segsize_means_10000, marker='o', markersize = 10, label='segsize_10000', color='mediumaquamarine', linewidth=5)


#REGRESSIONE LINEARE

# SEGSIZE 100
print('Segsize_100')
model_100 = LinearRegression()
# Addestrare il modello
model_100.fit(np.array(name_xaxis_100[1:]).reshape(-1,1), np.array(segsize_means_100[1:]).reshape(-1,1))
# Ottenere i coefficienti della regressione
slope_100 = model_100.coef_[0]
intercept_100 = model_100.intercept_
"""print(intercept_100)
print(slope_100)"""

# Calcolare i valori della linea di regressione
regression_line_100 = slope_100 * df_100['name_xaxis_100'] + intercept_100

# Visualizzare i risultati
"""plt.scatter((name_xaxis_100[1:]), (segsize_means_100[1:]), label='Dati di addestramento')
plt.plot((name_xaxis_100[1:]), regression_line_100, 'r-', label='Retaa di regressione 100', alpha = 0.3)"""



# SEGSIZE 1000
print('Segsize_1000')
model_1000 = LinearRegression()
# Addestrare il modello
model_1000.fit(np.array(name_xaxis_1000[1:]).reshape(-1,1), np.array(segsize_means_1000[1:]).reshape(-1,1))
# Ottenere i coefficienti della regressione
slope_1000 = model_1000.coef_[0]
intercept_1000 = model_1000.intercept_
"""print(intercept_1000)
print(slope_1000)"""

# Calcolare i valori della linea di regressione
regression_line_1000 = slope_1000 * df_1000['name_xaxis_1000'] + intercept_1000

# Visualizzare i risultati
"""plt.scatter((name_xaxis_1000[1:]), (segsize_means_1000[1:]), label='Dati di addestramento')
plt.plot((name_xaxis_1000[1:]), regression_line_1000, 'r-', label='Retaa di regressione 1000', alpha = 0.3)"""


# SEGSIZE 10000
print('Segsize_10000')
model_10000 = LinearRegression()
# Addestrare il modello
model_10000.fit(np.array(name_xaxis_10000[1:]).reshape(-1,1), np.array(segsize_means_10000[1:]).reshape(-1,1))
# Ottenere i coefficienti della regressione
slope_10000 = model_10000.coef_[0]
intercept_10000 = model_10000.intercept_
"""print(intercept_1000)
print(slope_1000)"""

# Calcolare i valori della linea di regressione
regression_line_10000 = slope_10000 * df_10000['name_xaxis_10000'] + intercept_10000

# Visualizzare i risultati
"""plt.scatter((name_xaxis_1000[1:]), (segsize_means_1000[1:]), label='Dati di addestramento')
plt.plot((name_xaxis_1000[1:]), regression_line_1000, 'r-', label='Retaa di regressione 1000', alpha = 0.3)"""




# FUNZIONE POLINOMIALE
# SEGSIZE 100
def log_poly(x, a, b, c):
    return a * np.log(b * x) + c

# Eseguire il fit della funzione polinomiale di grado logaritmico
params, covariance = curve_fit(log_poly, np.array(name_xaxis_100[1:]), np.array(segsize_means_100[1:]))

# Predire i valori con la funzione polinomiale di grado logaritmico
y_pred_100 = log_poly(np.array(name_xaxis_100[1:]), *params)

# plt.plot(name_xaxis_100[1:], y_pred_100, label='predict', color='lightblue', linewidth=3, linestyle='-', alpha=0.2)



# SEGSIZE 1000
def log_poly(x, a, b, c):
    return a * np.log(b * x) + c

# Eseguire il fit della funzione polinomiale di grado logaritmico
params, covariance = curve_fit(log_poly, np.array(name_xaxis_1000[1:]), np.array(segsize_means_1000[1:]))

# Predire i valori con la funzione polinomiale di grado logaritmico
y_pred_1000 = log_poly(np.array(name_xaxis_1000[1:]), *params)

# plt.plot(name_xaxis_1000[1:], y_pred_1000, label='predict', color='salmon', linewidth=3, linestyle='-', alpha=0.2)



# SEGSIZE 10000
def log_poly(x, a, b, c):
    return a * np.log(b * x) + c

# Eseguire il fit della funzione polinomiale di grado logaritmico
params, covariance = curve_fit(log_poly, np.array(name_xaxis_10000[1:]), np.array(segsize_means_10000[1:]))

# Predire i valori con la funzione polinomiale di grado logaritmico
y_pred_10000 = log_poly(np.array(name_xaxis_10000[1:]), *params)

# plt.plot(name_xaxis_10000[1:], y_pred_10000, label='predict', color='darkseagreen', linewidth=3, linestyle='-', alpha=0.2)



# Bontà di Adattamento

#R^2


#SEGSIZE 100
# FIT POLINOMIALE
print('SEGSIZE 100')
observed_data_100 = np.array(segsize_means_100[1:])
predicted_data_100 = y_pred_100
r2_100 = r2_score(observed_data_100, predicted_data_100)
print(r2_100)
#REGRESSIONE
predicted_data_100_reg = regression_line_100
r2_100_linear = r2_score(observed_data_100, predicted_data_100_reg)
print(r2_100_linear)


# FIT POLINOMIALE
#SEGSIZE 1000
print('SEGSIZE 1000')
observed_data_1000 = np.array(segsize_means_1000[1:])
predicted_data_1000 = y_pred_1000
r2_1000 = r2_score(observed_data_1000, predicted_data_1000)
print(r2_1000)
#REGRESSIONE
predicted_data_1000_reg = regression_line_1000
r2_1000_linear = r2_score(observed_data_1000, predicted_data_1000_reg)
print(r2_1000_linear)


# FIT POLINOMIALE
#SEGSIZE 10000
print('SEGSIZE 10000')
observed_data_10000 = np.array(segsize_means_10000[1:])
predicted_data_10000 = y_pred_10000
r2_10000 = r2_score(observed_data_10000, predicted_data_10000)
print(r2_10000)
#REGRESSIONE
predicted_data_10000_reg = regression_line_10000
r2_10000_linear = r2_score(observed_data_10000, predicted_data_10000_reg)
print(r2_10000_linear)



plt.xlabel('Cases', fontsize=24, labelpad=15)
plt.ylabel('Memory usage (MB)', fontsize=24, labelpad =15)
plt.yticks(fontsize=22)
plt.xticks(fontsize=22)
plt.legend(['Segment size 100 (KB)', 'Segment size 1000 (KB)', 'Segment size 10000 (KB)'], loc='upper left', fontsize=20)

plt.xlim([0, 9000.1])
plt.ylim([0, 70])

plt.tight_layout()
plt.grid(True, linestyle='--')

plt.savefig('scalability_numtraces.pdf')
plt.show()

