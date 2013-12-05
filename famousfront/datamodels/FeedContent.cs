using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedContent
  {
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
    public uint words
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint density
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint links
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public ulong status
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[] images
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia[] medias
    {
      get;set;
    }
  }
}