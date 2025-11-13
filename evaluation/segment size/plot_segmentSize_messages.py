import matplotlib.pyplot as plt


blue='deepskyblue'
red='tomato'
green='forestgreen'


# Sample data
x2 = ['0.5', '1','1.5', '2', '2.5', '3', '3.5', '4', '4.5', '5', '5.5', '6', '7']
log2 = [31, 24, 21, 20, 19, 19, 19, 19, 18, 18, 18, 18, 18]
bpic2 = [25, 21, 20, 19, 18, 18, 18, 18, 18, 18, 18, 18, 18]
sepsis2 = [24, 21, 20, 19, 19, 19, 18, 18, 18, 18, 18, 18, 18]

# Plotting
plt.style.use("seaborn-v0_8-bright")
plt.figure(figsize=(16, 9))
plt.grid(True, linestyle='--')


plt.plot(x2, log2, label='Motivating scenario', color=blue, linewidth=5, marker = 'p', markersize=10)
plt.plot(x2, bpic2, label='BPIC 2013', color=red, linewidth=5, marker = 'v', markersize=10)
plt.plot(x2, sepsis2, label='Sepsis', color=green, linewidth=5, marker = '*', markersize=10)
plt.xlabel('Segment size (MB)', fontsize = 30, labelpad= 15)
plt.ylabel('Number of exchanged messages', fontsize = 30, labelpad= 15)
plt.xticks(fontsize=25)
plt.yticks(fontsize=25)
plt.title('')
#The following command invert the Y axis
plt.legend(fontsize="25")  # Add legend to differentiate lines
#plt.grid(True)
plt.tight_layout()
plt.savefig('plot_segmentSize_messages.pdf')
plt.show()