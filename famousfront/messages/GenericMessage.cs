namespace famousfront.messages
{
  internal class GenericMessage
  {
    internal int code;
    internal string reason = null;
    protected GenericMessage()
    {
      code = 0;
    }
  }
}