db = open_db("file:memdb1?mode=memory&cache=shared")
fName = "./test_cases/01.Simple_examples.xlsx"
shName = "SimpleSheet"

q1 = db.load_excel_sheet(
    file_name = fName,
    sheet_name = shName,
    drop_table = True
)

# q1.print()

sha1 = q1.get_sha()

# print("Sha1:", sha1)

test_result = "sha of query:"+sha1

starlark_result_sha = sha1

# if sha1 == "ad211fb93f6948a9111606f217f178ad1f0cc0baf20a2ba3ad9e46d9e2a7a8ad":
#     test_result = "OK"

