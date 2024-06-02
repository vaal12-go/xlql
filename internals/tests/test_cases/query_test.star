
# newDB = open_db(file_name ="newDB.sqlite3")
newDB = open_db(file_name ="file:memdb1?mode=memory&cache=shared")
fName = "test_cases/01.Simple_examples.xlsx"
shName = "TableRange"
# print("Excel file name:", fName)
# print("\t has worksheets:", list_worksheets(fName))

# new_q = newDB.load_excel_sheet(
#     file_name=fName, 
#     sheet_name=shName,
#     skip_rows = 5
#     ) 
# print("1")
new_q = newDB.load_excel_sheet(
    file_name=fName, 
    sheet_name=shName,
    table_range = "C4:F8",
    drop_table = True
    ) 
# print(2)
# print("Query SHA:", new_q.get_sha())

starlark_result_sha = new_q.get_sha()


# print("starlark:", starlark_result_sha)