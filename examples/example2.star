print("This is an example file")

hrDB = open_db("hr_db2.sqlite")

departments_q = hrDB.load_excel_sheet(
    file_name = "Sample HR Database.xlsx",
    sheet_name = "departments"
)

employees_q = hrDB.load_excel_sheet(
    file_name = "Sample HR Database.xlsx",
    sheet_name = "employees"
)

combined_query = hrDB.run_query(
'''
SELECT 
    first_name,
    last_name,
    department_name
FROM
    employees
        INNER JOIN
    departments ON departments.department_id = employees.department_id
WHERE
    employees.department_id IN (1 , 2, 3);
'''
)

combined_query.save_to_excel("employee_dpt.xlsx", "employee_departments")



combined_query.print()