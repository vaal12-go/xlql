

test_result1 = TEST_getParametersLoad_excel_sheet(
    file_name = "qwe1.xlsx",
    sheet_name = "test1"
)

# print("SHA1:", test_result1)

test_result2 = TEST_getParametersLoad_excel_sheet(
    file_name = "qwe1.xlsx",
    sheet_name = "test1",
    skip_rows = 2
)

# print("SHA2:", test_result2)
  

test_result = "OK"