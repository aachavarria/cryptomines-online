- se puede tener mas de un warehouse, cuando deberia ser solo 1( revisa el gdd o dile al researcher que investigue)
- el UI al construir me dice que tengo 2 slots cuando realmente solo tengo 1
- puedo hacer investigaciones sin tener el technology center
- el chat no funciona. api.ts:436   GET http://localhost:5173/api/chat/messages?channel=world 500 (Internal Server Error), useChat.ts:27 
 Failed to fetch chat messages: AxiosError: Request failed with status code 500
    at async getChatMessages (api.ts:436:20)
    at async useChat.ts:18:22
    at async useChat.ts:49:9
    at async handleSend (ChatPanel.tsx:35:21)
- el UI del chat esta corrido hacia la izquierda, y ademas no se puede cerrar
- no puedo hacer claim de la primera quest. POST http://localhost:5173/api/quests/18088f24-78a5-4fc5-ab3b-bd4785354df9/claim 
- al construir una nave en la parte de modulos me salen en defense: orbital shield y atomic framework aun cuando no tenga ninguno de los 2
- al construir una nave en la parte de modulos me salen en auxiliary: nano station warehouse, station warehouse, super transmission engine aun cuando no tengo ninguno de los 3
- en military/blueprints donde veo todos mis blueprints no veo el nombre de cada uno 
- ademas me deja hacer research sin tener el weapon center( me parece que se necesita este revisa el gdd o dile al researcher que busque en wiki)
- comandantes dice coming soon pense que ya habiamos completado esto
- galaxy dice coming soon pense que ya habiamos completado esto
- en ship factory deberia poder ver una lista de las naves que tengo construidas y cuantas
- no puedo crear fleets. POST http://localhost:5173/api/fleets 400 (Bad Request) {"error":"invalid formation"}, para crear fleet segun entiendo hay q poner naves en uno de los 6 slots de la fleet pero no hay ui para eso.


nota: revisa la anterior, si no tenemos suficiente data hay q investigar para ver como funcionaba en galaxy online. en el gdd deberia estar casi todo



mensajes en el server:

[warning] build.bin is deprecated; set build.entrypoint instead

  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.64.5, built with Go go1.25.7

mkdir /home/yurei/cryptomines-online/backend/tmp
watching .
watching cmd
watching cmd/check_quests
watching cmd/server
watching internal
watching internal/combat
watching internal/database
watching internal/errs
watching internal/handlers
watching internal/middleware
watching internal/models
watching internal/services
watching internal/utils
watching internal/workers
!exclude tmp
building...
running...
2026/02/10 19:08:15 No .env file found, using environment variables
2026/02/10 19:08:15 Blueprint research worker started (checking every 30 seconds)
2026/02/10 19:08:15 Resource warehouse worker started (updating every 5 minutes)
2026/02/10 19:08:15 Server starting on :8080
2026/02/10 19:08:15 Updating warehouse for 1 resources
2026/02/10 19:08:15 Successfully updated warehouse for 1/1 resources
2026/02/10 19:13:14 Updating warehouse for 1 resources
2026/02/10 19:13:14 Successfully updated warehouse for 1/1 resources
2026/02/10 19:13:32 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:13:33 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:13:37 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:13:38 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:13:43 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:13:48 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:13:53 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:13:58 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:14:03 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:14:08 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:14:13 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:14:18 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:15:04 Failed to list available ships: pq: column ht.classification does not exist at position 5:4 (42703)
2026/02/10 19:15:04 Failed to list available ships: pq: column ht.classification does not exist at position 5:4 (42703)
2026/02/10 19:15:53 Quest 18088f24-78a5-4fc5-ab3b-bd4785354df9 completed for player 374f8321-c85a-47d7-bd8a-c7bd76b27097
2026/02/10 19:16:14 Failed to add item to inventory: pq: insert or update on table "player_inventory" violates foreign key constraint "player_inventory_item_key_fkey" (23503)
2026/02/10 19:16:14 Failed to unlock next quest: pq: current transaction is aborted, commands ignored until end of transaction block (25P02)
2026/02/10 19:16:14 Failed to commit: pq: could not complete operation in a failed transaction
2026/02/10 19:18:14 Updating warehouse for 1 resources
2026/02/10 19:18:14 Successfully updated warehouse for 1/1 resources
2026/02/10 19:23:13 Updating warehouse for 1 resources
2026/02/10 19:23:13 Successfully updated warehouse for 1/1 resources
2026/02/10 19:25:02 Failed to list available ships: pq: column ht.classification does not exist at position 5:4 (42703)
2026/02/10 19:25:02 Failed to list available ships: pq: column ht.classification does not exist at position 5:4 (42703)
2026/02/10 19:25:07 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:25:07 Failed to get chat messages: pq: syntax error at or near "$" at position 6:37 (42601)
2026/02/10 19:25:09 Failed to add item to inventory: pq: insert or update on table "player_inventory" violates foreign key constraint "player_inventory_item_key_fkey" (23503)
2026/02/10 19:25:09 Failed to unlock next quest: pq: current transaction is aborted, commands ignored until end of transaction block (25P02)
2026/02/10 19:25:09 Failed to commit: pq: could not complete operation in a failed transaction
