package models

func (t *TableModel) GetUpdatedFields(oldTable TableModel) (toAdd []Field, toUpdate [][2]Field, toDrop []Field) {
	// Create a map of the old fields by ID for faster lookup
	oldFieldMap := make(map[string]Field)
	for _, field := range oldTable.Fields {
		oldFieldMap[field.Id] = field
	}

	// Create a map to track IDs in the new table
	newFieldIds := make(map[string]bool)

	// Identify fields to add or update
	for _, newField := range t.Fields {
		newFieldIds[newField.Id] = true
		
		oldField, exists := oldFieldMap[newField.Id]
		if !exists {
			// Field doesn't exist in old table, so add it
			toAdd = append(toAdd, newField)
		} else if !oldField.Equals(&newField) {
			// Field exists but has changed, so update it
			toUpdate = append(toUpdate, [2]Field{oldField, newField})
		}
	}

	// Identify fields to drop (present in old table but not in new table)
	for _, oldField := range oldTable.Fields {
		if !newFieldIds[oldField.Id] {
			toDrop = append(toDrop, oldField)
		}
	}

	return toAdd, toUpdate, toDrop
}

func (f *Field) Equals(compareField *Field) bool {
    // Compare basic properties
    if f.Name != compareField.Name || f.Type != compareField.Type || f.Size != compareField.Size || 
       f.Nullable != compareField.Nullable || f.Primary != compareField.Primary || 
       f.Default != compareField.Default || f.Unique != compareField.Unique || 
       f.Index != compareField.Index || f.AutoIncrement != compareField.AutoIncrement {
        return false
    }
    
    // Compare foreign keys
    if (f.ForeignKey == nil) != (compareField.ForeignKey == nil) {
        return false
    }
    
    // Compare foreign keys cont.
    if f.ForeignKey != nil && compareField.ForeignKey != nil {
		if f.ForeignKey.Table != compareField.ForeignKey.Table || f.ForeignKey.Column != compareField.ForeignKey.Column {
			return false
		}
    }
    
    return true
}
