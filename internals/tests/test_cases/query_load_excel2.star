# print("Hello")


dbName = "tempDB.sqlite3"
newDB = open_db(dbName)
fName = "./test_cases/01.Simple_examples.xlsx"
shName = "SimpleSheet"


# print("Excel file name:", fName)
# print("\t has worksheets:", list_worksheets(fName))
# print("\n\n")
# print("DB has tables1:", newDB.get_tables())

new_q = newDB.load_excel_sheet(
    file_name=fName, 
    sheet_name=shName,
    drop_table = True
)

sha1 = new_q.get_sha()

# print("First SHA:", sha1)

new_q2 = newDB.run_query(
    """
        UPDATE SimpleSheet
            SET COL1 = "COL1_CHANGED"
    """
)


new_q = newDB.load_excel_sheet(
    file_name=fName, 
    sheet_name=shName,
    drop_table = True
)

sha2 = new_q.get_sha()

# print("Second SHA:", sha2)

# print("DB has tables2:", newDB.get_tables())

if (sha1==sha2):
    test_result = "OK"
else:
    test_result = "load_excel_sheet drop_table TEST FAILED"