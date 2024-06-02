db = open_db("file:memdb1?mode=memory&cache=shared")
fName = "./test_cases/01.Simple_examples.xlsx"
shName = "countries"

# print("Worksheets:", list_worksheets(fName))

q1 = db.load_excel_sheet(
    file_name = fName,
    sheet_name = shName,
    drop_table = True
)
sha1 = q1.get_sha()
# print("SHA1:", q1.get_sha())

# db.exec_sql()

db.exec_sql("""
    UPDATE countries
    SET REGION_ID = 5
    WHERE COUNTRY_NAME = "Argentina"
    """)

# q1.print()

sha2 = q1.get_sha()
# print("SHA2:", q1.get_sha())

if sha1 == "ca8b3013891ce3cfbd1a19e6b72b1269a8607249931b1004aea88e0a705d42eb" \
    and sha2 == "67947037bb80dfc0c8512ea6f624fbf8d085503d8ab80012003c67fcdf566777":
    test_result = "OK"
else :
    test_result = "sha1:"+sha1+"\nsha2:"+sha2