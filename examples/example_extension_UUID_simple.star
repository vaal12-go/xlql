#This requires UUID extension: https://github.com/nalgeon/sqlean/blob/main/docs/uuid.md
inMemDB = open_db("file:memdb1?mode=memory&cache=shared")

uuid_query = inMemDB.run_query(
    "select uuid_str(randomblob(16))"
)
print("UUID query:")
uuid_query.print()