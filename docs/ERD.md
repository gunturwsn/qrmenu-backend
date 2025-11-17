# QRMenu Backend – Entity Relationship Diagram

```mermaid
erDiagram
    TENANT ||--o{ TABLE : "has"
    TENANT ||--o{ CATEGORY : "organizes"
    TENANT ||--o{ ITEM : "owns"
    TENANT ||--o{ ADMIN_USER : "employs"
    TENANT ||--o{ ORDER : "receives"

    TABLE }o--|| TENANT : "belongs to"
    TABLE ||--o{ ORDER : "serves"

    CATEGORY }o--|| TENANT : "belongs to"
    CATEGORY ||--o{ ITEM : "groups"

    ITEM }o--|| TENANT : "belongs to"
    ITEM }o--|| CATEGORY : "classified under"
    ITEM ||--o{ ITEM_OPTION : "offers options"
    ITEM ||--o{ ORDER_ITEM : "included in"

    ITEM_OPTION }o--|| ITEM : "attached to"
    ITEM_OPTION ||--o{ ITEM_OPTION_VALUE : "provides values"

    ITEM_OPTION_VALUE }o--|| ITEM_OPTION : "belongs to"

    ORDER }o--|| TENANT : "belongs to"
    ORDER }o--|| TABLE : "placed from"
    ORDER ||--o{ ORDER_ITEM : "contains"

    ORDER_ITEM }o--|| ORDER : "part of"
    ORDER_ITEM }o--|| ITEM : "references"

    ADMIN_USER }o--|| TENANT : "assigned to"
```

## Entity Notes

- **Tenant**  
  Core partition key for the platform. Every other entity references a tenant to keep data isolated across restaurants/venues.

- **Table**  
  Physical table in a venue. Holds a unique token used by guests to fetch menus and place orders.

- **Category**  
  Groups menu items (e.g., Appetizers, Drinks). Each category belongs to a single tenant.

- **Item**  
  Actual menu entry. Belongs to both a tenant and a category. Can expose multiple options (sizes, add-ons).

- **ItemOption / ItemOptionValue**  
  Define configurable options for an item. An option belongs to an item and offers one or more values (e.g., `"Size" -> ["Small", "Large"]`).

- **Order / OrderItem**  
  Orders originate from a table and tenant. An order aggregates order items which point back to the item definition for pricing and naming.

- **AdminUser**  
  Staff member for a given tenant. Used for authentication and authorization across the admin endpoints.

This ERD mirrors the relationships encoded by the domain models inside `internal/domain`. Use it as a reference when extending repositories, adding migrations, or updating the OpenAPI specification.
