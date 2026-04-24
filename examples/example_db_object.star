#Example of Database object methods
# to run: 
# 

inMemDB = open_db(
    "file:memdb1?mode=memory&cache=shared"
)

mem_tbl1 = inMemDB.create_table(
    name = "SAMPLE_TABLE1",
    columns = {
	     "TXT_FIELD": "TEXT",
	     "NUM1": "NUMERIC",
	     "NUM2": "NUMERIC",
	     "NUM3": "NUMERIC"
	}
)

mem_tbl2 = inMemDB.create_table(
    name = "SAMPLE_TABLE2",
    columns = {
	     "TXT_FIELD": "TEXT",
	     "NUM1": "NUMERIC"
	}
)

for x in range(5):
    mem_tbl1.add_row("qwe1", x, x*10, x+1)

print("DB has following tables:", inMemDB.get_tables())

q1 = inMemDB.run_query(query = "SELECT * FROM SAMPLE_TABLE1")
q1.print()

tbl3 = inMemDB.create_table(
    name = "SMPL_TBL3",
    columns = {
        "TXT_COL1": "TEXT",
        "NUM_COL1": "NUMERIC"
    },
    auto_rename_table_name = False,
    drop_table = False
)

tbl3 = inMemDB.create_table(
    name = "SMPL_TBL3",
    columns = {
        "TXT_COL1": "TEXT",
        "NUM_COL1": "NUMERIC"
    },
    auto_rename_table_name = True,
    drop_table = False
)

print("DB has following tables:", inMemDB.get_tables())


tbl3 = inMemDB.create_table(
    name = "SMPL_TBL3",
    columns = {
        "TXT_COL1": "TEXT",
        "NUM_COL1": "NUMERIC",
        "NUM_COL2": "NUMERIC"
    },
    auto_rename_table_name = False,
    drop_table = True
)

print("DB has following tables:", inMemDB.get_tables())