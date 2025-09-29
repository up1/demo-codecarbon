import os, time
import numpy as np, pandas as pd

N = 300_000
rng = np.random.default_rng(0)
df = pd.DataFrame({
    "id": np.arange(N, dtype=np.int64),
    "int_col": rng.integers(0, 1_000_000, size=N, dtype=np.int32),
    "float_col": rng.random(N) * 1e6,
    "bool_col": rng.choice([True, False], size=N),
    "cat_col": pd.Categorical(rng.integers(0, 50, size=N)),
    "str_col": rng.choice([f"s{i}" for i in range(2000)], size=N),
    "date_col": pd.to_datetime("2022-01-01") + pd.to_timedelta(rng.integers(0, 365, size=N), unit="D"),
})
df["ts_col"] = df["date_col"] + pd.to_timedelta(rng.integers(0, 86_400, size=N), unit="s")

csv_path = "data.csv"
pq_path  = "data.parquet"

def t(f,*a,**k): s=time.perf_counter(); r=f(*a,**k); return r, time.perf_counter()-s

# write
_, t_csv_w = t(df.to_csv, csv_path, index=False)
_, t_pq_w  = t(df.to_parquet, pq_path, engine="pyarrow", compression="snappy", index=False)

# read full
_, t_csv_r = t(pd.read_csv, csv_path)
_, t_pq_r  = t(pd.read_parquet, pq_path, engine="pyarrow")

# read only two columns
cols = ["id","float_col"]
_, t_csv_cols = t(pd.read_csv, csv_path, usecols=cols)
_, t_pq_cols  = t(pd.read_parquet, pq_path, columns=cols, engine="pyarrow")

# sizes
csv_mb = os.path.getsize(csv_path)/1024/1024
pq_mb  = os.path.getsize(pq_path)/1024/1024

# Create a DataFrame with the results
results = pd.DataFrame({
    'Operation': ['Write', 'Read Full', 'Read Columns'],
    'CSV (seconds)': [round(t_csv_w, 3), round(t_csv_r, 3), round(t_csv_cols, 3)],
    'Parquet (seconds)': [round(t_pq_w, 3), round(t_pq_r, 3), round(t_pq_cols, 3)]
})

# Add file sizes as a separate row
size_row = pd.DataFrame({
    'Operation': ['File Size (MB)'],
    'CSV (seconds)': [round(csv_mb, 2)],
    'Parquet (seconds)': [round(pq_mb, 2)]
})

results = pd.concat([results, size_row])

# Display the table
print(results.to_string(index=False))