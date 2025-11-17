add new command:

    app read

first screen, list available translations:

    1. VI1934
    2. NKJV
    3. NIV

This list is just the list of folders older appConfig.listeners.store.filesystem.directory

When user enter 1 or 1 or 3, user will enter the second screen, list of books

    1. Sáng Thế Ký
    2. Xuất Ê-díp-tô ký
    3. …

This is list of books found under translations directory. Book name is parsed by `_category_.json > label` under the book's directory. 
