import { Container, Paper, Typography } from "@mui/material";

interface ErrorViewProps {
  msg: string;
}

function ErrorView({ msg }: ErrorViewProps) {
  return (
    <Container maxWidth="lg">
      <Paper elevation={12} sx={{ padding: 2, mt: 2, borderRadius: 2 }}>
        <Typography component="h1" variant="h3" sx={{ mb: 2 }}>
          Error
        </Typography>
        <Typography component="p" variant="body1">
          {msg}
        </Typography>
      </Paper>
    </Container>
  );
}

export default ErrorView;
