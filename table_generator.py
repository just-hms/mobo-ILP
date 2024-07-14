import pandas as pd

# Read the data into a DataFrame
df = pd.read_csv('data/benchmark.csv')

keyGroups = [
    ['nOuts'], 
    ['nIn'], 
    ['onesRatio'],
]

for g in keyGroups:
    # Group by the specified columns and calculate the mean
    grouped_df = df.groupby(g).agg({'multiCost': 'mean', 'singleCost': 'mean', 'duration[ms]' : 'mean'}).reset_index()

    grouped_df['improvement[\\%]'] = ((grouped_df['singleCost'] - grouped_df['multiCost']) / grouped_df['singleCost']) * 100

    grouped_df = grouped_df.drop(columns=['singleCost', 'multiCost'])

    latex_table = grouped_df.to_latex(index=False, float_format="%.2f")

    print(latex_table)

    print('\\vspace{0.2cm}')
    print()


# Calculate the global mean without grouping
global_mean_df = df[['multiCost', 'singleCost']].mean().to_frame().T
global_mean_df['improvement[\\%]'] = ((global_mean_df['singleCost'] - global_mean_df['multiCost']) / global_mean_df['singleCost']) * 100

global_mean_table = global_mean_df.to_latex(index=False, float_format="%.2f")

print(global_mean_table)
