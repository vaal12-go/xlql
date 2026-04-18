inMemDB = open_db(
    "file:memdb1?mode=memory&cache=shared"
)

mem_q = inMemDB.run_query("SELECT sqlite_version()")
mem_q.print()