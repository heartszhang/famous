using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedMedia
  {
    [DataMember(EmitDefaultValue = false)]
    public string title
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string description
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string uri
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string local
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int width
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int height
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int duration
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string mime
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string thumbanil
    {
      get;set;
    }
  }
}