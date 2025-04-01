## Code Optimization Checklist

1. Put structure into the code to follow clean code and clean architecture -cohesion and coupling
3. What design patterns can we introduce to improve the codebase readability an extensibility 
4. Configuration management across different deployment environments 
5. What optimisations can benefits this code to serve 1 million concurrent users/sec - both from code side and infrastructure side.
6. Is this code observable enough for production environment - Struggled logging with easily configurable log level across different environment, and can different different logging agent pick the logs up easily,
7. Security considerations
8. Metrics and server endpoint should be separated
9. Adpat to 12 factor app - https://12factor.net/
