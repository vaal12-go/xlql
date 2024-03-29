print("Hello2 ")


dbName = "tempDB.sqlite3"
newDB = open_db(dbName)
fName = "./test_cases/01.Simple_examples.xlsx"
shName = "SimpleSheet"


print("Excel file name:", fName)
print("\t has worksheets:", list_worksheets(fName))
print("\n\n")
print("DB has tables1:", newDB.get_tables())

new_q = newDB.load_excel_sheet(
    file_name=fName, 
    sheet_name=shName,
    drop_table = True
)

sha1 = new_q.get_sha()

print("First SHA:", sha1)

new_new_q = newDB.create_table(
    name = shName,
    columns = {
        "NewCol1": "TEXT",
        "NewCol23": "NUMERIC",
        "NewCol2": "NUMERIC",
        "NewCol0": "NUMERIC"
	},
    drop_table = True
)

new_new_q.add_row("New data in Col1", 1, 2, 3)

print("DB has tables2:", newDB.get_tables())

sha2 = new_new_q.get_sha()

print(" new_new_q Second SHA:", sha2)

new3_q = newDB.run_query("SELECT * FROM "+shName)

sha3 = new3_q.get_sha()

print("Have sha#3:", sha3)

print("DB has tables2:", newDB.get_tables())

if (sha3=="fe56525ffcd9dbf6d6d317bead6d8cb15c9d2a781f862acadaa0d672c6b656e2"):
    test_result = "OK"
else:
    test_result = "create_table drop_table TEST FAILED"

