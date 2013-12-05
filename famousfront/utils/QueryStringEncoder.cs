namespace famousfront.utils
{
  internal static class QueryStringEncoder
  {
    internal static string QueryString(this object request)
    {
      return HttpQueryParameter.EncodeFieldQueryString(request);
    }
    internal static string QueryPath(this object request, string pattren)
    {
      return HttpQueryParameter.EncodeRelativePath(request, pattren);
    }
  }
}