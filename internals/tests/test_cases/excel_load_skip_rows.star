db = open_db("file:memdb1?mode=memory&cache=shared")
fName = "./test_cases/01.Simple_examples.xlsx"
shName = "SkipRows"

q1 = db.load_excel_sheet(
    file_name = fName,
    sheet_name = shName,
    skip_rows = 4,
    drop_table = True
)

# q1.print()

starlark_result_sha = q1.get_sha()
 

# print("starlark_result_sha:", starlark_result_sha)

test_result = starlark_result_sha

#
# if starlark_result_sha == "6aa5aaf2e4eb1a575e81c89153119622a54b0e681eea2d4947ef3f275491ab18":
#     test_result = "OK"

